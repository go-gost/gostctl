package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateRecorder(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriRecorder
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

func (c *Client) UpdateRecorder(ctx context.Context, recorder string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriRecorder + "/" + recorder
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

func (c *Client) DeleteRecorder(ctx context.Context, recorder string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriRecorder + "/" + recorder
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
