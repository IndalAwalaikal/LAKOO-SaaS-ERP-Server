package dto

import "time"

type InventoryMutationRequest struct {
	ProductID    string  `json:"product_id" binding:"required"`
	MutationType string  `json:"mutation_type" binding:"required,oneof=in out adjustment"`
	Qty          float64 `json:"qty" binding:"required,gt=0"`
	Reference    string  `json:"reference"`
	Notes        string  `json:"notes"`
}

type InventoryMutationResponse struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	MutationType string    `json:"mutation_type"`
	Qty          float64   `json:"qty"`
	Balance      float64   `json:"balance"`
	Reference    string    `json:"reference"`
	Notes        string    `json:"notes"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}
