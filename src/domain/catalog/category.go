package catalog

import "time"

type Category struct {
	ID          string  // UUID
	Name        string
	Description string
	ParentID    *string // UUID for self-reference (Hierarchical 1, 2, 3)
	Level       int     // 1, 2, 3
	ImageURL    string
	Active      bool
	CreatedBy   int     // FK to User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SearchResultCategory struct {
	Data       *[]Category
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
