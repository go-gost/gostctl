package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateChain(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriChain
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

func (c *Client) UpdateChain(ctx context.Context, chain string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriChain + "/" + chain
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

func (c *Client) DeleteChain(ctx context.Context, chain string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriChain + "/" + chain
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
