package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"lakoo/backend/internal/domain"
)

type SaleRepository interface {
	CreateWithItems(sale *domain.Sale, items []domain.SaleItem) error
	FindByTenant(tenantID string) ([]domain.Sale, error)
	GetRecent(tenantID string, limit int) ([]domain.Sale, error)
}

type saleRepository struct {
	db *sqlx.DB
}

func NewSaleRepository(db *sqlx.DB) SaleRepository {
	return &saleRepository{db: db}
}

func (r *saleRepository) CreateWithItems(sale *domain.Sale, items []domain.SaleItem) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// Make sure to rollback if err occurs, else commit.
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	querySale := `
		INSERT INTO sales (id, tenant_id, invoice_no, cashier_id, customer_id, total_amount, discount_amount, grand_total, payment_method, payment_status, created_at, updated_at)
		VALUES (:id, :tenant_id, :invoice_no, :cashier_id, :customer_id, :total_amount, :discount_amount, :grand_total, :payment_method, :payment_status, :created_at, :updated_at)
	`
	_, err = tx.NamedExec(querySale, sale)
	if err != nil {
		return fmt.Errorf("failed to insert sale: %w", err)
	}

	queryItem := `
		INSERT INTO sale_items (id, sale_id, product_id, qty, price, subtotal, created_at)
		VALUES (:id, :sale_id, :product_id, :qty, :price, :subtotal, :created_at)
	`
	for _, item := range items {
		_, err = tx.NamedExec(queryItem, &item)
		if err != nil {
			return fmt.Errorf("failed to insert sale item: %w", err)
		}

		// Update product stock
		queryStock := `UPDATE products SET stock_qty = stock_qty - ? WHERE id = ? AND tenant_id = ?`
		_, err = tx.Exec(queryStock, item.Qty, item.ProductID, sale.TenantID)
		if err != nil {
			return fmt.Errorf("failed to reduce product stock: %w", err)
		}
	}

	return nil
}

func (r *saleRepository) FindByTenant(tenantID string) ([]domain.Sale, error) {
	var sales []domain.Sale
	err := r.db.Select(&sales, "SELECT * FROM sales WHERE tenant_id = ? AND deleted_at IS NULL ORDER BY created_at DESC", tenantID)
	if err != nil {
		return nil, err
	}
	if sales == nil {
		sales = []domain.Sale{}
	}
	return sales, nil
}

func (r *saleRepository) GetRecent(tenantID string, limit int) ([]domain.Sale, error) {
	var sales []domain.Sale
	query := "SELECT * FROM sales WHERE tenant_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ?"
	err := r.db.Select(&sales, query, tenantID, limit)
	if err != nil {
		return nil, err
	}
	if sales == nil {
		sales = []domain.Sale{}
	}
	return sales, nil
}
