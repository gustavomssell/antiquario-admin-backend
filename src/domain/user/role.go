package user

import "time"

type Role struct {
	ID          string // UUID
	Name        string // e.g., admin, vendedor, comprador
	Description string
	Permissions string // Stored as JSON string representation
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
