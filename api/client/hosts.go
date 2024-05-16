package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateHostMapper(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriHosts
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	return nil
}

func (c *Client) UpdateHostMapper(ctx context.Context, hostMapper string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriHosts + "/" + hostMapper
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	return nil
}

func (c *Client) DeleteHostMapper(ctx context.Context, hostMapper string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriHosts + "/" + hostMapper
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	return nil
}
