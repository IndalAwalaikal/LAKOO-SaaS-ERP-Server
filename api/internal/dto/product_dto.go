package dto

type ProductRequest struct {
	SKU          string  `json:"sku"`
	Barcode      string  `json:"barcode"`
	ImageURL     *string `json:"image_url"`
	Name         string  `json:"name" binding:"required"`
	CostPrice    float64 `json:"cost_price" binding:"required,min=0"`
	SellingPrice float64 `json:"selling_price" binding:"required,min=0"`
	StockQty     float64 `json:"stock_qty"`
	MinStock     float64 `json:"min_stock"`
	Unit         string  `json:"unit" binding:"required"`
	IsActive     *bool   `json:"is_active"`
}
