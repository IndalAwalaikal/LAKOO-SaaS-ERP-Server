package domain

import "time"

type InventoryMutation struct {
	ID           string    `json:"id" db:"id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	ProductID    string    `json:"product_id" db:"product_id"`
	MutationType string    `json:"mutation_type" db:"mutation_type"` // in, out, adjustment, sale, return
	Qty          float64   `json:"qty" db:"qty"`
	Balance      float64   `json:"balance" db:"balance"` // stock balance AFTER mutation
	Reference    string    `json:"reference" db:"reference"` // invoice number or manually entered ref
	Notes        string    `json:"notes" db:"notes"`
	CreatedBy    string    `json:"created_by" db:"created_by"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
