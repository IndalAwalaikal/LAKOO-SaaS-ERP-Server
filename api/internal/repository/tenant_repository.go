package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"lakoo/backend/internal/domain"
)

type TenantRepository interface {
	Create(tenant *domain.Tenant) error
	GetByID(id string) (*domain.Tenant, error)
	GetBySlug(slug string) (*domain.Tenant, error)
	Update(tenant *domain.Tenant) error
}

type tenantRepository struct {
	db *sqlx.DB
}

func NewTenantRepository(db *sqlx.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(t *domain.Tenant) error {
	query := `
		INSERT INTO tenants (id, slug, name, plan, status, owner_id, payment_config, logo_url, trial_ends_at, created_at, updated_at) 
		VALUES (:id, :slug, :name, :plan, :status, :owner_id, :payment_config, :logo_url, :trial_ends_at, :created_at, :updated_at)`
	
	_, err := r.db.NamedExec(query, t)
	return err
}

func (r *tenantRepository) GetByID(id string) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.db.Get(&t, "SELECT * FROM tenants WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Or a specific ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *tenantRepository) GetBySlug(slug string) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.db.Get(&t, "SELECT * FROM tenants WHERE slug = ? AND deleted_at IS NULL", slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *tenantRepository) Update(t *domain.Tenant) error {
	query := `UPDATE tenants SET name = :name, slug = :slug, payment_config = :payment_config, logo_url = :logo_url, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExec(query, t)
	return err
}
