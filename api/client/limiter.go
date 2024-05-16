package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateLimiter(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriLimiter
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

func (c *Client) UpdateLimiter(ctx context.Context, limiter string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriLimiter + "/" + limiter
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

func (c *Client) DeleteLimiter(ctx context.Context, limiter string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriLimiter + "/" + limiter
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
