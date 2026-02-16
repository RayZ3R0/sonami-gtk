package v1

type Album struct {
	Cover string `json:"cover"`
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type AlbumItem struct {
	Album          Album    `json:"album"`
	AllowStreaming bool     `json:"allowStreaming"`
	Artists        []Artist `json:"artists"`
	Duration       int      `json:"duration"`
	ID             int      `json:"id"`
	Title          string   `json:"title"`
	TrackNumber    int      `json:"trackNumber"`
}
