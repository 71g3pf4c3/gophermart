package entity

// Balance represents the current and withdrawn bonus amounts.
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
