package openapi

import "codeberg.org/dergs/tidalwave/pkg/tidalapi/helper"

type Playlist struct {
	Data     TypedObject[PlaylistAttributes] `json:"data"`
	Included []Object
	Links    map[string]string `json:"links"`
}

type PlaylistAttributes struct {
	CreatedAt         helper.OpenAPIDateTime `json:"createdAt"`
	Description       string                 `json:"description"`
	Duration          helper.DurationISO8601 `json:"duration"`
	Name              string                 `json:"name"`
	NumberOfFollowers int                    `json:"numberOfFollowers"`
	NumberOfItems     int                    `json:"numberOfItems"`
}
