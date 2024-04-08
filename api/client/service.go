package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateService(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriService
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

func (c *Client) UpdateService(ctx context.Context, service string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriService + "/" + service
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

func (c *Client) DeleteService(ctx context.Context, service string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriService + "/" + service
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
