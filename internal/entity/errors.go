package entity

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderOwnedBySelf   = errors.New("order already uploaded by this user")
	ErrOrderOwnedByOther  = errors.New("order already uploaded by another user")
	ErrInvalidOrderNumber = errors.New("invalid order number")

	ErrInsufficientFunds = errors.New("insufficient funds")
)
