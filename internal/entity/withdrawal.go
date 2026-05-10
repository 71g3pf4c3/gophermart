package entity

import "time"

// Withdrawal represents a balance withdrawal operation.
type Withdrawal struct {
	ID          string    `json:"-"`
	UserID      string    `json:"-"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
