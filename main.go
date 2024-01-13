package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	_ "net"
	"os"
	"time"

	"gioui.org/app"
	_ "gioui.org/app/permission/storage"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/go-gost/gui/api/client"
	"github.com/go-gost/gui/api/runner"
	"github.com/go-gost/gui/config"
	"github.com/go-gost/gui/ui"
)

func main() {
	Init()

	go func() {
		w := app.NewWindow(
			app.Title("GOST"),
			app.MinSize(800, 600),
		)
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	go func() {
		for e := range runner.Default().Event() {
			if e.Err == nil {
				w.Invalidate()
			} else {
				slog.Error(fmt.Sprintf("task: %s", e.Err), "task", e.TaskID)
			}
		}
	}()

	ui := ui.NewUI()
	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func Init() {
	config.Init()

	cfg := config.Global()
	if cfg.CurrentServer >= 0 && cfg.CurrentServer < len(cfg.Servers) {
		server := cfg.Servers[cfg.CurrentServer]
		client.SetDefault(client.NewClient(server.URL, server.Timeout))
		interval := server.Interval
		if interval <= 0 {
			interval = 3 * time.Second
		}
		runner.Default().ExecAsync(context.Background(), runner.GetConfigTask("getconfig"), interval)
	}
}
