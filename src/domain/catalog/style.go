package catalog

import "time"

type Style struct {
	ID            string // UUID
	Name          string // e.g., Colonial Brasileiro, Francês
	Description   string
	PeriodID      *string // FK to Period
	OriginCountry string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SearchResultStyle struct {
	Data       *[]Style
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
