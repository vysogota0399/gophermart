package models

type Order struct {
	Number     string  `json:"number"`
	State      string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at,omitempty"`
}
