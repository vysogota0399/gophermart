package models

type Withdraw struct {
	OrderNumber string  `json:"order_number,omitempty"`
	Sum         float64 `json:"sum,omitempty"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}
