package product

import "time"

type ProductMedia struct {
	ID            string // UUID
	ProductID     string // FK to Product
	MediaType     string // e.g., image, video
	MediaURL      string
	AltText       string
	IsPrimary     bool
	OrderPosition int
	FileSize      int64
	Duration      *int // Duration in seconds (for videos)
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
