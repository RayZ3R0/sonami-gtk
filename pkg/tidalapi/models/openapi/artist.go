package openapi

type Artist struct {
	Data     TypedObject[ArtistAttributes] `json:"data"`
	Included []Object
	Links    map[string]string `json:"links"`
}

type ArtistAttributes struct {
	Name       string  `json:"name"`
	Popularity float64 `json:"popularity"`
}
