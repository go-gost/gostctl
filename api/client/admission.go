package client

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) CreateAdmission(ctx context.Context, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriAdmission
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

func (c *Client) UpdateAdmission(ctx context.Context, admission string, body io.Reader) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriAdmission + "/" + admission
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

func (c *Client) DeleteAdmission(ctx context.Context, admission string) error {
	if c.url == "" {
		return nil
	}

	url := c.url + uriAdmission + "/" + admission
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
