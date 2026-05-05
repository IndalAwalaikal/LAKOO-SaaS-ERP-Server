package dto

type CustomerRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"omitempty,email"`
	Phone   string `json:"phone" binding:"omitempty"`
	Address string `json:"address" binding:"omitempty"`
	IsMember bool   `json:"is_member"`
}

type CustomerResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Points  int    `json:"points"`
	IsMember bool  `json:"is_member"`
}
