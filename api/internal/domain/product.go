package domain

import (
	"time"
)

type Product struct {
	ID           string     `json:"id" db:"id"`
	TenantID     string     `json:"tenant_id" db:"tenant_id"`
	SKU          string     `json:"sku" db:"sku"`
	Barcode      string     `json:"barcode" db:"barcode"`
	ImageURL     *string    `json:"image_url" db:"image_url"`
	Name         string     `json:"name" db:"name"`
	CostPrice    float64    `json:"cost_price" db:"cost_price"`
	SellingPrice float64    `json:"selling_price" db:"selling_price"`
	StockQty     float64    `json:"stock_qty" db:"stock_qty"`
	MinStock     float64    `json:"min_stock" db:"min_stock"`
	Unit         string     `json:"unit" db:"unit"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}
