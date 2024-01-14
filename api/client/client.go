package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
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

type Options struct {
	Userinfo *url.Userinfo
	Timeout  time.Duration
}

type Option func(opts *Options)

func WithUserinfo(userinfo *url.Userinfo) Option {
	return func(opts *Options) {
		opts.Userinfo = userinfo
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

type Client struct {
	client   http.Client
	url      string
	userinfo *url.Userinfo
}

func NewClient(url string, opts ...Option) *Client {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	if url != "" {
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
	}

	return &Client{
		client: http.Client{
			Timeout: options.Timeout,
		},
		url:      url,
		userinfo: options.Userinfo,
	}
}

func (c *Client) GetConfig(ctx context.Context) (*api.Config, error) {
	if c.url == "" {
		return &api.Config{}, nil
	}

	url := c.url + "config"
	// slog.Debug(fmt.Sprintf("GET %s", url))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if c.userinfo != nil {
		username := c.userinfo.Username()
		password, _ := c.userinfo.Password()
		req.SetBasicAuth(username, password)
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
