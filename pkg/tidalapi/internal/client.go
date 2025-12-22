package internal

import "net/http"

type Client struct {
	*http.Client
	countryCode string
	token       string
}

func NewClient(countryCode string, token string) *Client {
	client := &http.Client{
		Transport: MiddlewareRoundTripper{countryCode: countryCode},
	}

	return &Client{
		Client:      client,
		countryCode: countryCode,
		token:       token,
	}
}
