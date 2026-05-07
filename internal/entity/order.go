package entity

import "time"

// OrderStatus is the current accrual processing state.
type OrderStatus string

const (
	// OrderStatusNew means the order has been accepted but not processed yet.
	OrderStatusNew OrderStatus = "NEW"
	// OrderStatusProcessing means the order is being processed by the accrual service.
	OrderStatusProcessing OrderStatus = "PROCESSING"
	// OrderStatusInvalid means the accrual service rejected the order.
	OrderStatusInvalid OrderStatus = "INVALID"
	// OrderStatusProcessed means the accrual service successfully processed the order.
	OrderStatusProcessed OrderStatus = "PROCESSED"
)

// Order represents a user's uploaded order.
type Order struct {
	Number     string      `json:"number"`
	UserID     string      `json:"user_id,omitempty"`
	Status     OrderStatus `json:"status"`
	Accrual    *float64    `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}
