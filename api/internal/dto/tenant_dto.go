package dto

type RegisterRequest struct {
	TenantName string `json:"tenant_name" binding:"required"`
	Slug       string `json:"slug" binding:"required"`
	OwnerName  string `json:"owner_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
}

type RegisterResponse struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
	Slug     string `json:"slug"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token  string      `json:"token"`
	User   UserObj     `json:"user"`
	Tenant TenantObj   `json:"tenant"`
}

type UserObj struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type TenantObj struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	PaymentConfig *string `json:"payment_config"`
	LogoURL       *string `json:"logo_url"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UpdateTenantRequest struct {
	Name          string  `json:"name" binding:"required"`
	Slug          string  `json:"slug" binding:"required"`
	OwnerName     string  `json:"owner_name" binding:"required"`
	Email         string  `json:"email" binding:"required,email"`
	PaymentConfig string  `json:"payment_config"`
	LogoURL       string  `json:"logo_url"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type AddStaffRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required"` // manager or cashier
}

type StaffMemberResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type SalesTrendResponse struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}
