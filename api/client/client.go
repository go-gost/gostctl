package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-gost/gui/api"
)

var (
	defaultClient atomic.Value
)

func init() {
	defaultClient.Store(&Client{})
}

func SetDefault(c *Client) {
	defaultClient.Store(c)
}

func Default() *Client {
	v, _ := defaultClient.Load().(*Client)
	return v
}

type Client struct {
	client http.Client
	url    string
}

func NewClient(url string, timeout time.Duration) *Client {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return &Client{
		client: http.Client{
			Timeout: timeout,
		},
		url: url,
	}
}

func (c *Client) GetConfig(ctx context.Context) (*api.Config, error) {
	if c.url == "" {
		return nil, nil
	}

	url := c.url + "config"
	slog.Debug(fmt.Sprintf("GET %s", url))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	cfg := &api.Config{}
	if err := json.NewDecoder(resp.Body).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
