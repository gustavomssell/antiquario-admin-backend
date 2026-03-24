package product

import "time"

type ProductDocument struct {
	ID           string // UUID
	ProductID    string // FK to Product
	DocumentType string // e.g., certificate, appraisal, invoice, provenance
	DocumentURL  string
	Title        string
	Description  string
	FileSize     int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
