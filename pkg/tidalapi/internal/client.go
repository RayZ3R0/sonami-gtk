package internal

import (
	"net/http"

	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/auth"
)

type Client struct {
	*http.Client
	authStrategies []auth.AuthStrategy
	countryCode    string
}

func NewClient(countryCode string, authStrategies ...auth.AuthStrategy) *Client {
	client := &http.Client{
		Transport: MiddlewareRoundTripper{authStrategies: append(authStrategies, auth.NewOpenAPIClientAuthStrategy()), countryCode: countryCode},
	}

	return &Client{
		Client:         client,
		authStrategies: authStrategies,
		countryCode:    countryCode,
	}
}
