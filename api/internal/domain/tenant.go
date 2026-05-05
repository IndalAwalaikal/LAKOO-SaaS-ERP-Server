package domain

import (
	"time"
)

type Tenant struct {
	ID          string     `json:"id" db:"id"`
	Slug        string     `json:"slug" db:"slug"`
	Name        string     `json:"name" db:"name"`
	Plan        string     `json:"plan" db:"plan"`
	Status      string     `json:"status" db:"status"`
	OwnerID       string     `json:"owner_id" db:"owner_id"`
	PaymentConfig *string    `json:"payment_config" db:"payment_config"`
	LogoURL       *string    `json:"logo_url" db:"logo_url"`
	TrialEndsAt   *time.Time `json:"trial_ends_at" db:"trial_ends_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}

type User struct {
	ID        string     `json:"id" db:"id"`
	TenantID  string     `json:"tenant_id" db:"tenant_id"`
	Name      string     `json:"name" db:"name"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"-" db:"password"`
	Role      string     `json:"role" db:"role"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
