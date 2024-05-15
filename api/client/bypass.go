package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateBypass(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriBypass
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

func (c *Client) UpdateBypass(ctx context.Context, bypass string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriBypass + "/" + bypass
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

func (c *Client) DeleteBypass(ctx context.Context, bypass string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriBypass + "/" + bypass
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
