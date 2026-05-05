package domain

import "time"

type FinanceTransaction struct {
	ID          string     `json:"id" db:"id"`
	TenantID    string     `json:"tenant_id" db:"tenant_id"`
	Type        string     `json:"type" db:"type"`             // income, expense
	Category    string     `json:"category" db:"category"`     // e.g. operational, salary, rent, marketing
	Amount      float64    `json:"amount" db:"amount"`
	Description string     `json:"description" db:"description"`
	Date        time.Time  `json:"date" db:"date"`
	ReferenceID string     `json:"reference_id" db:"reference_id"` // links to a sale ID or payroll ID
	CreatedBy   string     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}
