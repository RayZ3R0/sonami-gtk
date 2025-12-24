package openapi

type Track struct {
	Data     TypedObject[TrackAttributes] `json:"data"`
	Included []Object
	Links    map[string]string `json:"links"`
}

type TrackAttributes struct {
	Duration  string   `json:"duration"`
	Explicit  bool     `json:"explicit"`
	MediaTags []string `json:"media_tags"`
	Title     string   `json:"title"`
}
