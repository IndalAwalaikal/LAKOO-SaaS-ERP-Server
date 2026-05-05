package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"lakoo/backend/internal/domain"
	"lakoo/backend/internal/dto"
	"lakoo/backend/internal/repository"
)

type ProductService interface {
	CreateProduct(tenantID string, req *dto.ProductRequest) (*domain.Product, error)
	UpdateProduct(id, tenantID string, req *dto.ProductRequest) (*domain.Product, error)
	DeleteProduct(id, tenantID string) error
	GetProducts(tenantID string) ([]domain.Product, error)
	GetProductByID(id, tenantID string) (*domain.Product, error)
}

type productService struct {
	repo        repository.ProductRepository
	financeRepo repository.FinanceRepository
}

func NewProductService(repo repository.ProductRepository, fr repository.FinanceRepository) ProductService {
	return &productService{repo: repo, financeRepo: fr}
}

func (u *productService) CreateProduct(tenantID string, req *dto.ProductRequest) (*domain.Product, error) {
	now := time.Now()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	product := &domain.Product{
		ID:           uuid.New().String(),
		TenantID:     tenantID,
		SKU:          req.SKU,
		Barcode:      req.Barcode,
		ImageURL:     req.ImageURL,
		Name:         req.Name,
		CostPrice:    req.CostPrice,
		SellingPrice: req.SellingPrice,
		StockQty:     req.StockQty,
		MinStock:     req.MinStock,
		Unit:         req.Unit,
		IsActive:     isActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := u.repo.Create(product)
	if err != nil {
		return nil, err
	}

	// Record automated expense if stock > 0
	if product.StockQty > 0 {
		exp := &domain.FinanceTransaction{
			ID:          uuid.New().String(),
			TenantID:    tenantID,
			Type:        "Expense",
			Category:    "Pembelian Stok (Otomatis)",
			Amount:      product.StockQty * product.CostPrice,
			Description: "Pengeluaran stok awal saat pendaftaran produk: " + product.Name,
			Date:        now,
			ReferenceID: product.ID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		_ = u.financeRepo.Create(exp)
	}

	return product, nil
}

func (u *productService) UpdateProduct(id, tenantID string, req *dto.ProductRequest) (*domain.Product, error) {
	existing, err := u.repo.FindByID(id, tenantID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("product not found")
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	existing.SKU = req.SKU
	existing.Barcode = req.Barcode
	existing.ImageURL = req.ImageURL
	existing.Name = req.Name
	existing.CostPrice = req.CostPrice
	existing.SellingPrice = req.SellingPrice
	existing.StockQty = req.StockQty
	existing.MinStock = req.MinStock
	existing.Unit = req.Unit
	existing.IsActive = isActive
	existing.UpdatedAt = time.Now()

	diffStock := req.StockQty - existing.StockQty

	err = u.repo.Update(existing)
	if err != nil {
		return nil, err
	}

	// Record automated expense if stock increased
	if diffStock > 0 {
		now := time.Now()
		exp := &domain.FinanceTransaction{
			ID:          uuid.New().String(),
			TenantID:    tenantID,
			Type:        "Expense",
			Category:    "Pembelian Stok (Otomatis)",
			Amount:      diffStock * req.CostPrice,
			Description: "Penambahan stok produk via update: " + existing.Name,
			Date:        now,
			ReferenceID: existing.ID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		_ = u.financeRepo.Create(exp)
	}

	return existing, nil
}

func (u *productService) DeleteProduct(id, tenantID string) error {
	return u.repo.Delete(id, tenantID)
}

func (u *productService) GetProducts(tenantID string) ([]domain.Product, error) {
	return u.repo.FindByTenant(tenantID)
}

func (u *productService) GetProductByID(id, tenantID string) (*domain.Product, error) {
	product, err := u.repo.FindByID(id, tenantID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}
