package service

import (
	"time"

	"github.com/google/uuid"
	"lakoo/backend/internal/domain"
	"lakoo/backend/internal/dto"
	"lakoo/backend/internal/repository"
)

type InventoryService interface {
	AdjustStock(tenantID, userID string, req *dto.InventoryMutationRequest) (*domain.InventoryMutation, error)
	GetProductHistory(tenantID, productID string) ([]domain.InventoryMutation, error)
}

type inventoryService struct {
	repo        repository.InventoryRepository
	productRepo repository.ProductRepository
	financeRepo repository.FinanceRepository
}

func NewInventoryService(repo repository.InventoryRepository, pr repository.ProductRepository, fr repository.FinanceRepository) InventoryService {
	return &inventoryService{repo: repo, productRepo: pr, financeRepo: fr}
}

func (u *inventoryService) AdjustStock(tenantID, userID string, req *dto.InventoryMutationRequest) (*domain.InventoryMutation, error) {
	mut := &domain.InventoryMutation{
		ID:           uuid.New().String(),
		TenantID:     tenantID,
		ProductID:    req.ProductID,
		MutationType: req.MutationType,
		Qty:          req.Qty,
		Reference:    req.Reference,
		Notes:        req.Notes,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
	}

	err := u.repo.AdjustStockTx(mut)
	if err != nil {
		return nil, err
	}

	// Record automated expense if addition
	if req.MutationType == "in" {
		product, err := u.productRepo.FindByID(req.ProductID, tenantID)
		if err == nil && product != nil {
			exp := &domain.FinanceTransaction{
				ID:          uuid.New().String(),
				TenantID:    tenantID,
				Type:        "Expense",
				Category:    "Penambahan Stok (Otomatis)",
				Amount:      req.Qty * product.CostPrice,
				Description: "Pembelian stok baru via manajemen inventori: " + product.Name,
				Date:        time.Now(),
				ReferenceID: mut.ID,
				CreatedBy:   userID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			_ = u.financeRepo.Create(exp)
		}
	}

	return mut, nil
}

func (u *inventoryService) GetProductHistory(tenantID, productID string) ([]domain.InventoryMutation, error) {
	return u.repo.GetProductMutations(tenantID, productID)
}
