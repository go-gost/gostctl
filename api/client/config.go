package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

func (c *Client) SaveConfig(ctx context.Context, filepath string) error {
	if c.url == "" {
		return nil
	}

	format := "yaml"
	if strings.HasSuffix(filepath, ".json") {
		format = "json"
	}
	values := url.Values{}
	if format != "" {
		values.Set("format", format)
	}
	if filepath != "" {
		values.Set("path", filepath)
	}

	surl := c.url + uriConfig
	if len(values) > 0 {
		surl += ("?" + values.Encode())
	}

	// slog.Debug(fmt.Sprintf("GET %s", url))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, surl, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	rsp := api.Response{}
	if err := json.NewDecoder(resp).Decode(&rsp); err != nil {
		return err
	}
	if rsp.Code != 0 {
		return fmt.Errorf("%d %v", rsp.Code, rsp.Msg)
	}

	return nil
}
