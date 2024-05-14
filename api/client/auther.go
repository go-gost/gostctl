package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateAuther(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriAuther
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

func (c *Client) UpdateAuther(ctx context.Context, auther string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriAuther + "/" + auther
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

func (c *Client) DeleteAuther(ctx context.Context, auther string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriAuther + "/" + auther
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
