package v2

type PageItem struct {
	Icons         []any    `json:"icons"`
	Items         []Item   `json:"items"`
	ModuleID      string   `json:"moduleId"`
	ShowQuickPlay bool     `json:"showQuickPlay,omitempty"`
	Subtitle      string   `json:"subtitle,omitempty"`
	Title         string   `json:"title"`
	Type          ItemType `json:"type"`
	ViewAll       string   `json:"viewAll"`
}
