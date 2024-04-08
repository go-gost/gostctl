package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-gost/gostctl/api"
)

func (c *Client) GetConfig(ctx context.Context) (*api.Config, error) {
	if c.url == "" {
		return &api.Config{}, nil
	}

	url := c.url + uriConfig
	// slog.Debug(fmt.Sprintf("GET %s", url))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	cfg := &api.Config{}
	if err := json.NewDecoder(resp).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
