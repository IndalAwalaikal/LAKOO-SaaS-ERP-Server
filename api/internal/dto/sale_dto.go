package dto

type SaleItemRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Qty       float64 `json:"qty" binding:"required,min=0.01"`
}

type SaleRequest struct {
	CustomerID     *string           `json:"customer_id"`
	DiscountAmount float64           `json:"discount_amount"`
	PaymentMethod  string            `json:"payment_method" binding:"required,oneof=cash qris transfer ewallet"`
	PaymentStatus  string            `json:"payment_status" binding:"omitempty,oneof=paid pending refunded"`
	Items          []SaleItemRequest `json:"items" binding:"required,min=1"`
}
