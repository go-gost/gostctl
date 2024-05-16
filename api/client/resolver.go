package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateResolver(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriResolver
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

func (c *Client) UpdateResolver(ctx context.Context, resolver string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriResolver + "/" + resolver
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

func (c *Client) DeleteResolver(ctx context.Context, resolver string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriResolver + "/" + resolver
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
