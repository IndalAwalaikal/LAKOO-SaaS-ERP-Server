package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"lakoo/backend/internal/domain"
)

type InventoryRepository interface {
	InsertMutation(tx *sqlx.Tx, mut *domain.InventoryMutation) error
	GetProductMutations(tenantID, productID string) ([]domain.InventoryMutation, error)
	AdjustStockTx(mut *domain.InventoryMutation) error
}

type inventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) InsertMutation(tx *sqlx.Tx, mut *domain.InventoryMutation) error {
	query := `
		INSERT INTO inventory_mutations (id, tenant_id, product_id, mutation_type, qty, balance, reference, notes, created_by, created_at)
		VALUES (:id, :tenant_id, :product_id, :mutation_type, :qty, :balance, :reference, :notes, :created_by, :created_at)
	`
	if tx != nil {
		_, err := tx.NamedExec(query, mut)
		return err
	}
	_, err := r.db.NamedExec(query, mut)
	return err
}

func (r *inventoryRepository) GetProductMutations(tenantID, productID string) ([]domain.InventoryMutation, error) {
	var mutations []domain.InventoryMutation
	query := "SELECT * FROM inventory_mutations WHERE tenant_id = ? AND product_id = ? ORDER BY created_at DESC"
	err := r.db.Select(&mutations, query, tenantID, productID)
	
	if mutations == nil {
		mutations = []domain.InventoryMutation{}
	}
	return mutations, err
}

// AdjustStockTx handles a manual inventory adjustment in a transaction
func (r *inventoryRepository) AdjustStockTx(mut *domain.InventoryMutation) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	
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

	// 1. Lock and read product current stock
	var currentStock float64
	err = tx.Get(&currentStock, "SELECT stock_qty FROM products WHERE id = ? AND tenant_id = ? FOR UPDATE", mut.ProductID, mut.TenantID)
	if err != nil {
		return fmt.Errorf("could not retrieve product stock: %v", err)
	}

	// 2. Adjust current stock
	if mut.MutationType == "in" || mut.MutationType == "adjustment" {
		currentStock += mut.Qty
	} else if mut.MutationType == "out" {
		currentStock -= mut.Qty
		if currentStock < 0 {
			return fmt.Errorf("insufficient stock to perform this transaction")
		}
	}

	mut.Balance = currentStock

	// 3. Update the product
	_, err = tx.Exec("UPDATE products SET stock_qty = ?, updated_at = ? WHERE id = ? AND tenant_id = ?", 
		currentStock, mut.CreatedAt, mut.ProductID, mut.TenantID)
	if err != nil {
		return fmt.Errorf("failed to update product stock: %v", err)
	}

	// 4. Save mutation log
	err = r.InsertMutation(tx, mut)
	if err != nil {
		return fmt.Errorf("failed to save inventory log: %v", err)
	}

	return nil
}
