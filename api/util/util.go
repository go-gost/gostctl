package util

import (
	"context"
	"net/url"
	"time"

	"github.com/go-gost/gui/api"
	"github.com/go-gost/gui/api/client"
	"github.com/go-gost/gui/api/runner"
	"github.com/go-gost/gui/config"
)

func RestartGetConfigTask() {
	api.SetConfig(&api.Config{})

	var server config.Server
	cfg := config.Global()
	if cfg.CurrentServer >= 0 && cfg.CurrentServer < len(cfg.Servers) {
		server = cfg.Servers[cfg.CurrentServer]
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
	runner.Default().Exec(context.Background(), runner.GetConfigTask(),
		runner.WithAync(true),
		runner.WithInterval(interval),
	)
}
