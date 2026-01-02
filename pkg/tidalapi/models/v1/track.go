package v1

type Track struct {
	Album struct {
		Cover string `json:"cover"`
		ID    int    `json:"id"`
		Title string `json:"title"`
	}
	Artists []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	Duration int    `json:"duration"`
	ID       int    `json:"id"`
	Title    string `json:"title"`
}

type TrackMix struct {
	ID string `json:"id"`
}
