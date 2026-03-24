package user

import "time"

type UserProfile struct {
	ID             string // UUID
	UserID         int    // FK to user (int)
	Phone          string // JSON string representation
	Address        string // JSON string representation
	DocumentType   string
	DocumentNumber string
	BirthDate      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
