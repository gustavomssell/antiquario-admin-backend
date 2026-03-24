package product

import "time"

type ProductSet struct {
	ID                string // UUID
	SetCode           string // e.g., ANT-2024-SET-001
	Name              string
	Description       string
	TotalPieces       int
	CanSellSeparately bool
	CreatedBy         int // FK to User
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
