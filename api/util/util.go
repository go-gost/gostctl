package util

import (
	"context"
	"net/url"
	"time"

	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/api/client"
	"github.com/go-gost/gostctl/api/runner"
	"github.com/go-gost/gostctl/api/runner/task"
	"github.com/go-gost/gostctl/config"
)

func RestartGetConfigTask() {
	api.SetConfig(&api.Config{})

	var server *config.Server
	cfg := config.Get()
	if cfg.CurrentServer >= 0 && cfg.CurrentServer < len(cfg.Servers) {
		server = cfg.Servers[cfg.CurrentServer]
	}
	if server == nil {
		return
	}

	var userinfo *url.Userinfo
	if server.Username != "" {
		userinfo = url.UserPassword(server.Username, server.Password)
	}
	client.SetDefault(client.NewClient(server.URL,
		client.WithTimeout(server.Timeout),
		client.WithUserinfo(userinfo),
	))
	interval := server.Interval
	if interval <= 0 {
		interval = 3 * time.Second
	}
	runner.Exec(context.Background(), task.GetConfig(),
		runner.WithAync(true),
		runner.WithInterval(interval),
		runner.WithCancel(true),
	)

	if server.AutoSave != "" {
		runner.Exec(context.Background(),
			task.SaveConfig(server.AutoSave),
			runner.WithAync(true),
			runner.WithCancel(true),
		)
	}
}
