package search_results

import "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"

type SearchResults struct {
	client *internal.Client
}

func New(client *internal.Client) *SearchResults {
	return &SearchResults{
		client: client,
	}
}
