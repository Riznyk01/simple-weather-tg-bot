package http_client

import (
	"net/http"
	"time"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type DefaultHTTPClient struct {
	Timeout time.Duration
}

func (c *DefaultHTTPClient) Get(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: c.Timeout,
	}
	return client.Get(url)
}
