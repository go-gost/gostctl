package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateHop(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriHop
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

func (c *Client) UpdateHop(ctx context.Context, hop string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriHop + "/" + hop
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

func (c *Client) DeleteHop(ctx context.Context, hop string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriHop + "/" + hop
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
