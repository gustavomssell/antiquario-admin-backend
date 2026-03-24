package catalog

import "time"

type Period struct {
	ID          string // UUID
	Name        string // e.g., Barroco, Art Déco
	Description string
	StartYear   *int
	EndYear     *int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SearchResultPeriod struct {
	Data       *[]Period
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
