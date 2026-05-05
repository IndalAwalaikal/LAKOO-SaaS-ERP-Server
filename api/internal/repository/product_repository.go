package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"lakoo/backend/internal/domain"
)

type ProductRepository interface {
	Create(p *domain.Product) error
	Update(p *domain.Product) error
	Delete(id, tenantID string) error
	FindByID(id, tenantID string) (*domain.Product, error)
	FindByTenant(tenantID string) ([]domain.Product, error)
	FindLowStock(tenantID string) ([]domain.Product, error)
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(p *domain.Product) error {
	query := `
		INSERT INTO products (id, tenant_id, sku, barcode, image_url, name, cost_price, selling_price, stock_qty, min_stock, unit, is_active, created_at, updated_at)
		VALUES (:id, :tenant_id, :sku, :barcode, :image_url, :name, :cost_price, :selling_price, :stock_qty, :min_stock, :unit, :is_active, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(query, p)
	return err
}

func (r *productRepository) Update(p *domain.Product) error {
	query := `
		UPDATE products
		SET sku = :sku, barcode = :barcode, image_url = :image_url, name = :name, cost_price = :cost_price, 
		    selling_price = :selling_price, stock_qty = :stock_qty, min_stock = :min_stock, 
		    unit = :unit, is_active = :is_active, updated_at = :updated_at
		WHERE id = :id AND tenant_id = :tenant_id AND deleted_at IS NULL
	`
	result, err := r.db.NamedExec(query, p)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("product not found or unauthorized")
	}
	return nil
}

func (r *productRepository) Delete(id, tenantID string) error {
	query := `
		UPDATE products SET deleted_at = ?
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
		return errors.New("product not found or unauthorized")
	}
	return nil
}

func (r *productRepository) FindByID(id, tenantID string) (*domain.Product, error) {
	var p domain.Product
	err := r.db.Get(&p, "SELECT * FROM products WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL", id, tenantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // soft missing
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) FindByTenant(tenantID string) ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Select(&products, "SELECT * FROM products WHERE tenant_id = ? AND deleted_at IS NULL ORDER BY name ASC", tenantID)
	if err != nil {
		return nil, err
	}
	// Return empty slice instead of nil for better JSON marshaling
	if products == nil {
		products = []domain.Product{}
	}
	return products, nil
}

func (r *productRepository) FindLowStock(tenantID string) ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Select(&products, "SELECT * FROM products WHERE tenant_id = ? AND stock_qty <= min_stock AND deleted_at IS NULL ORDER BY stock_qty ASC", tenantID)
	if err != nil {
		return nil, err
	}
	if products == nil {
		products = []domain.Product{}
	}
	return products, nil
}
