package dto

import "time"

type FinanceRequest struct {
	Type        string  `json:"type" binding:"required,oneof=income expense"`
	Category    string  `json:"category" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required"`
	Date        string  `json:"date" binding:"required"` // ISO string format expecting YYYY-MM-DD
	ReferenceID string  `json:"reference_id" binding:"omitempty"`
}

type FinanceResponse struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}
