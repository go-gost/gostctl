package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-gost/gostctl/api"
)

const (
	uriConfig    = "/config"
	uriService   = uriConfig + "/services"
	uriChain     = uriConfig + "/chains"
	uriHop       = uriConfig + "/hops"
	uriAuther    = uriConfig + "/authers"
	uriAdmission = uriConfig + "/admissions"
	uriBypass    = uriConfig + "/bypasses"
	uriResolver  = uriConfig + "/resolvers"
	uriHosts     = uriConfig + "/hosts"
	uriLimiter   = uriConfig + "/limiters"
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
		url = strings.TrimSuffix(url, "/")
	}

	return &Client{
		client: http.Client{
			Timeout: options.Timeout,
		},
		url:      url,
		userinfo: options.Userinfo,
	}
}

func (c *Client) do(req *http.Request) (io.ReadCloser, error) {
	if req == nil {
		return nil, nil
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

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		var rsp api.Response
		json.NewDecoder(resp.Body).Decode(&rsp)
		return nil, fmt.Errorf("%d %s", rsp.Code, rsp.Msg)
	}

	return resp.Body, nil
}
