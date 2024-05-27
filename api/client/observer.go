package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateObserver(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriObserver
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

func (c *Client) UpdateObserver(ctx context.Context, observer string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriObserver + "/" + observer
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

func (c *Client) DeleteObserver(ctx context.Context, observer string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriObserver + "/" + observer
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
