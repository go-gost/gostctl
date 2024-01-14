package config

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gioui.org/app"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

const (
	configFile = "gost.yml"
	logFile    = "gost.log"
)

var (
	configDir string
)

func Init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))

	dir, err := app.DataDir()
	if err != nil {
		slog.Error(fmt.Sprintf("appDir: %v", err))
	}
	if dir == "" {
		dir, _ = os.Getwd()
	}
	configDir = filepath.Join(dir, "gost")
	os.MkdirAll(configDir, 0755)

	slog.Info(fmt.Sprintf("appDir: %s", configDir))

	if err := global.load(); err != nil {
		slog.Error(fmt.Sprintf("load config: %v", err))
		if _, ok := err.(*os.PathError); ok {
			global.Write()
		}
	}

	initLog()
}

func initLog() {
	cfg := global.Log
	if cfg == nil {
		return
	}

	/*
		logDir := filepath.Join(configDir, "logs")
		os.MkdirAll(logDir, 0755)
		slog.Info(fmt.Sprintf("log dir: %s", logDir))
	*/

	var out io.Writer
	switch cfg.Output {
	case "none", "null":
		out = io.Discard
	case "stdout", "":
		out = os.Stdout
	case "stderr":
		out = os.Stderr
	default:
		if cfg.Rotation != nil {
			out = &lumberjack.Logger{
				Filename:   cfg.Output,
				MaxSize:    cfg.Rotation.MaxSize,
				MaxAge:     cfg.Rotation.MaxAge,
				MaxBackups: cfg.Rotation.MaxBackups,
				LocalTime:  cfg.Rotation.LocalTime,
				Compress:   cfg.Rotation.Compress,
			}
		} else {
			os.MkdirAll(filepath.Dir(cfg.Output), 0755)
			f, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				slog.Warn(fmt.Sprintf("open log file %s: %v", cfg.Output, err))
			} else {
				out = f
			}
		}
	}

	level := &slog.LevelVar{}
	switch cfg.Level {
	case "debug":
		level.Set(slog.LevelDebug)
	case "warn":
		level.Set(slog.LevelWarn)
	case "error":
		level.Set(slog.LevelError)
	case "info":
		fallthrough
	default:
		level.Set(slog.LevelInfo)
	}

	var handler slog.Handler
	if cfg.Format == "text" {
		handler = slog.NewTextHandler(out, &slog.HandlerOptions{AddSource: true, Level: level})
	} else {
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{AddSource: true, Level: level})
	}

	slog.SetDefault(slog.New(handler))
}

var (
	global    = &Config{}
	globalMux sync.RWMutex
)

func Global() *Config {
	globalMux.RLock()
	defer globalMux.RUnlock()

	cfg := &Config{}
	*cfg = *global
	return cfg
}

func Set(c *Config) {
	globalMux.Lock()
	defer globalMux.Unlock()

	global = c
}

type Settings struct {
	Lang  string
	Theme string
}

type Server struct {
	Name     string
	URL      string        `yaml:"url"`
	Username string        `yaml:",omitempty"`
	Password string        `yaml:",omitempty"`
	Interval time.Duration `yaml:",omitempty"`
	Timeout  time.Duration `yaml:",omitempty"`
}

type Log struct {
	Output   string
	Level    string
	Format   string
	Rotation *LogRotation
}

type LogRotation struct {
	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `yaml:"maxSize,omitempty" json:"maxSize,omitempty"`
	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `yaml:"maxAge,omitempty" json:"maxAge,omitempty"`
	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `yaml:"maxBackups,omitempty" json:"maxBackups,omitempty"`
	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time. The default is to use UTC
	// time.
	LocalTime bool `yaml:"localTime,omitempty" json:"localTime,omitempty"`
	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `yaml:"compress,omitempty" json:"compress,omitempty"`
}

type Config struct {
	Servers       []Server
	CurrentServer int `yaml:"currentServer"`
	Settings      *Settings
	Log           *Log
}

func (c *Config) load() error {
	f, err := os.Open(filepath.Join(configDir, configFile))
	if err != nil {
		return err
	}
	defer f.Close()

	return yaml.NewDecoder(f).Decode(c)
}

func (c *Config) Write() error {
	if c == nil {
		c = &Config{}
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	defer enc.Close()

	enc.SetIndent(2)
	if err := enc.Encode(c); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(configDir, configFile), buf.Bytes(), 0644)
}
