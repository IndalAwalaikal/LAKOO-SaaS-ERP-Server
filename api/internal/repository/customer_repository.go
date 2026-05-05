package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"lakoo/backend/internal/domain"
)

type CustomerRepository interface {
	Create(c *domain.Customer) error
	Update(c *domain.Customer) error
	Delete(id, tenantID string) error
	FindByID(id, tenantID string) (*domain.Customer, error)
	FindByTenant(tenantID string) ([]domain.Customer, error)
}

type customerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(c *domain.Customer) error {
	query := `
		INSERT INTO customers (id, tenant_id, name, email, phone, address, points, is_member, created_at, updated_at)
		VALUES (:id, :tenant_id, :name, :email, :phone, :address, :points, :is_member, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(query, c)
	return err
}

func (r *customerRepository) Update(c *domain.Customer) error {
	query := `
		UPDATE customers 
		SET name = :name, email = :email, phone = :phone, address = :address, points = :points, is_member = :is_member, updated_at = :updated_at
		WHERE id = :id AND tenant_id = :tenant_id AND deleted_at IS NULL
	`
	result, err := r.db.NamedExec(query, c)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("customer not found or unauthorized")
	}
	return nil
}

func (r *customerRepository) Delete(id, tenantID string) error {
	query := `
		UPDATE customers SET deleted_at = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, time.Now(), id, tenantID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("customer not found or unauthorized")
	}
	return nil
}

func (r *customerRepository) FindByID(id, tenantID string) (*domain.Customer, error) {
	var c domain.Customer
	err := r.db.Get(&c, "SELECT * FROM customers WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL", id, tenantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *customerRepository) FindByTenant(tenantID string) ([]domain.Customer, error) {
	var customers []domain.Customer
	err := r.db.Select(&customers, "SELECT * FROM customers WHERE tenant_id = ? AND deleted_at IS NULL ORDER BY name ASC", tenantID)
	if err != nil {
		return nil, err
	}
	if customers == nil {
		customers = []domain.Customer{}
	}
	return customers, nil
}
