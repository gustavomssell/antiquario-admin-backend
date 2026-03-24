package catalog

import "time"

type Material struct {
	ID          string // UUID
	Name        string // e.g., Madeira de Lei, Bronze
	Description string
	Category    string // e.g., madeira, metal, tecido
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SearchResultMaterial struct {
	Data       *[]Material
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
