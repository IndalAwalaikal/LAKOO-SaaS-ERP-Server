package dto

import "time"

type NotificationResponse struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // low_stock, new_sale
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	IsRead    bool      `json:"is_read"`
}
