package v2

type Page struct {
	UUID string `json:"uuid"`
	Page struct {
		Cursor string `json:"cursor,omitempty"`
	} `json:"page"`
	Items []PageItem `json:"items"`
}
