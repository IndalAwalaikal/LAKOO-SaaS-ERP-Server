package domain

import (
	"time"
)

type Sale struct {
	ID             string     `json:"id" db:"id"`
	TenantID       string     `json:"tenant_id" db:"tenant_id"`
	InvoiceNo      string     `json:"invoice_no" db:"invoice_no"`
	CashierID      string     `json:"cashier_id" db:"cashier_id"`
	CustomerID     *string    `json:"customer_id" db:"customer_id"`
	TotalAmount    float64    `json:"total_amount" db:"total_amount"`
	DiscountAmount float64    `json:"discount_amount" db:"discount_amount"`
	GrandTotal     float64    `json:"grand_total" db:"grand_total"`
	PaymentMethod  string     `json:"payment_method" db:"payment_method"`
	PaymentStatus  string     `json:"payment_status" db:"payment_status"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" db:"deleted_at"`
}

type SaleItem struct {
	ID        string    `json:"id" db:"id"`
	SaleID    string    `json:"sale_id" db:"sale_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	Qty       float64   `json:"qty" db:"qty"`
	Price     float64   `json:"price" db:"price"`
	Subtotal  float64   `json:"subtotal" db:"subtotal"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
