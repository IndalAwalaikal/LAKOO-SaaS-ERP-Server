package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"lakoo/backend/internal/domain"
)

type FinanceRepository interface {
	Create(ft *domain.FinanceTransaction) error
	Delete(id, tenantID string) error
	FindByTenant(tenantID string) ([]domain.FinanceTransaction, error)
	GetMetrics(tenantID string, startDate, endDate time.Time) (income, expense float64, err error)
}

type financeRepository struct {
	db *sqlx.DB
}

func NewFinanceRepository(db *sqlx.DB) FinanceRepository {
	return &financeRepository{db: db}
}

func (r *financeRepository) Create(ft *domain.FinanceTransaction) error {
	query := `
		INSERT INTO finance_transactions (id, tenant_id, type, category, amount, description, date, reference_id, created_by, created_at, updated_at)
		VALUES (:id, :tenant_id, :type, :category, :amount, :description, :date, :reference_id, :created_by, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(query, ft)
	return err
}

func (r *financeRepository) Delete(id, tenantID string) error {
	query := `UPDATE finance_transactions SET deleted_at = ? WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`
	res, err := r.db.Exec(query, time.Now(), id, tenantID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("transaction not found")
	}
	return nil
}

func (r *financeRepository) FindByTenant(tenantID string) ([]domain.FinanceTransaction, error) {
	var txs []domain.FinanceTransaction
	err := r.db.Select(&txs, "SELECT * FROM finance_transactions WHERE tenant_id = ? AND deleted_at IS NULL ORDER BY date DESC, created_at DESC", tenantID)
	if txs == nil {
		txs = []domain.FinanceTransaction{}
	}
	return txs, err
}

func (r *financeRepository) GetMetrics(tenantID string, startDate, endDate time.Time) (income float64, expense float64, err error) {
	query := `
		SELECT type, SUM(amount) as total 
		FROM finance_transactions 
		WHERE tenant_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL 
		GROUP BY type
	`
	
	type Result struct {
		Type  string  `db:"type"`
		Total float64 `db:"total"`
	}
	
	var res []Result
	err = r.db.Select(&res, query, tenantID, startDate, endDate)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, 0, err
	}

	for _, r := range res {
		if r.Type == "Income" {
			income = r.Total
		} else if r.Type == "Expense" {
			expense = r.Total
		}
	}
	return income, expense, nil
}
