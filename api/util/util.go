package util

import (
	"context"
	"time"

	"github.com/go-gost/gui/api/client"
	"github.com/go-gost/gui/api/runner"
	"github.com/go-gost/gui/config"
)

func RestartGetConfigTask() {
	cfg := config.Global()
	if cfg.CurrentServer >= 0 && cfg.CurrentServer < len(cfg.Servers) {
		server := cfg.Servers[cfg.CurrentServer]
		client.SetDefault(client.NewClient(server.URL, server.Timeout))
		interval := server.Interval
		if interval <= 0 {
			interval = 3 * time.Second
		}
		runner.Default().Exec(context.Background(), runner.GetConfigTask(),
			runner.WithAync(true),
			runner.WithInterval(interval),
		)
	}
}
