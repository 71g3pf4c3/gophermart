package entity

import "time"

// User represents an account in the loyalty system.
type User struct {
	ID           string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username     string    `json:"username" example:"johndoe"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at" example:"2026-01-01T00:00:00Z"`
}
