package mpsanity

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	ProjectID string
	Dataset   string
	Token     string

	HTTPClient *http.Client
}

type Option interface {
	Apply(c *Client) error
}

type WithDataset string

func (d WithDataset) Apply(c *Client) error {
	c.Dataset = string(d)
	return nil
}

type WithToken string

func (t WithToken) Apply(c *Client) error {
	c.Token = string(t)
	return nil
}

func New(projectID string, opts ...Option) (*Client, error) {
	c := &Client{
		ProjectID:  projectID,
		HTTPClient: &http.Client{},
	}

	for _, o := range opts {
		if err := o.Apply(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) newRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	u := fmt.Sprintf("https://%s.api.sanity.io/v1%s", c.ProjectID, path)
	r, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Accept", "application/json")
	if c.Token != "" {
		r.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return r, nil
}
