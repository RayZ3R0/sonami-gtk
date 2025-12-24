package openapi

import "codeberg.org/dergs/tidalwave/pkg/tidalapi/helper"

type Album struct {
	Data     TypedObject[AlbumAttributes] `json:"data"`
	Included []Object
	Links    map[string]string `json:"links"`
}

type AlbumAttributes struct {
	Duration        string              `json:"duration"`
	Explicit        bool                `json:"explicit"`
	MediaTags       []string            `json:"media_tags"`
	NumberOfItems   int                 `json:"numberOfItems"`
	NumberOfVolumes int                 `json:"numberOfVolumes"`
	Popularity      float64             `json:"popularity"`
	ReleaseDate     helper.TimeDateOnly `json:"releaseDate"`
	Title           string              `json:"title"`
	Type            string              `json:"type"`
}
