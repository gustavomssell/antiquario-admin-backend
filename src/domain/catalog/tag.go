package catalog

import "time"

type Tag struct {
	ID          string // UUID
	Name        string // e.g., raro, restaurado
	Color       string // Hex color for UI
	Description string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SearchResultTag struct {
	Data       *[]Tag
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
