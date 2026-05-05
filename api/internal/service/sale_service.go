package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"lakoo/backend/internal/domain"
	"lakoo/backend/internal/dto"
	"lakoo/backend/internal/repository"
)

type SaleService interface {
	CreateSale(tenantID string, cashierID string, req *dto.SaleRequest) (*domain.Sale, error)
	ListSales(tenantID string) ([]domain.Sale, error)
	GetSalesTrend(tenantID string) ([]dto.SalesTrendResponse, error)
}

type saleService struct {
	saleRepo    repository.SaleRepository
	productRepo repository.ProductRepository
	financeRepo repository.FinanceRepository
}

func NewSaleService(sr repository.SaleRepository, pr repository.ProductRepository, fr repository.FinanceRepository) SaleService {
	return &saleService{
		saleRepo:    sr,
		productRepo: pr,
		financeRepo: fr,
	}
}

func (u *saleService) CreateSale(tenantID string, cashierID string, req *dto.SaleRequest) (*domain.Sale, error) {
	saleID := uuid.New().String()
	now := time.Now()
	invoiceNo := fmt.Sprintf("INV-%s-%d", tenantID[:8], now.Unix())

	var items []domain.SaleItem
	var totalAmount float64
	var totalCOGS float64

	for _, itemReq := range req.Items {
		product, err := u.productRepo.FindByID(itemReq.ProductID, tenantID)
		if err != nil || product == nil {
			return nil, fmt.Errorf("invalid product: %s", itemReq.ProductID)
		}

		if product.StockQty < itemReq.Qty {
			return nil, fmt.Errorf("insufficient stock for product: %s", product.Name)
		}

		subtotal := product.SellingPrice * itemReq.Qty
		totalAmount += subtotal
		totalCOGS += product.CostPrice * itemReq.Qty

		item := domain.SaleItem{
			ID:        uuid.New().String(),
			SaleID:    saleID,
			ProductID: product.ID,
			Qty:       itemReq.Qty,
			Price:     product.SellingPrice,
			Subtotal:  subtotal,
			CreatedAt: now,
		}
		items = append(items, item)
	}

	grandTotal := totalAmount - req.DiscountAmount
	if grandTotal < 0 {
		return nil, errors.New("discount cannot exceed total amount")
	}

	status := "paid"
	if req.PaymentStatus != "" {
		status = req.PaymentStatus
	}

	sale := &domain.Sale{
		ID:             saleID,
		TenantID:       tenantID,
		InvoiceNo:      invoiceNo,
		CashierID:      cashierID,
		CustomerID:     req.CustomerID,
		TotalAmount:    totalAmount,
		DiscountAmount: req.DiscountAmount,
		GrandTotal:     grandTotal,
		PaymentMethod:  req.PaymentMethod,
		PaymentStatus:  status,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := u.saleRepo.CreateWithItems(sale, items)
	if err != nil {
		return nil, err
	}

	// 1. Record Income
	incomeRecord := &domain.FinanceTransaction{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Type:        "Income",
		Category:    "Penjualan Kasir",
		Amount:      grandTotal,
		Description: fmt.Sprintf("Pendapatan otomatis (Metode: %s) dari Nomor Struk: %s", req.PaymentMethod, invoiceNo),
		Date:        now,
		ReferenceID: saleID,
		CreatedBy:   cashierID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = u.financeRepo.Create(incomeRecord)
	if err != nil {
		fmt.Printf("Warning: failed to automatically log income for sale %s: %v\n", saleID, err)
	}

	return sale, nil
}

func (u *saleService) ListSales(tenantID string) ([]domain.Sale, error) {
	return u.saleRepo.FindByTenant(tenantID)
}

func (u *saleService) GetSalesTrend(tenantID string) ([]dto.SalesTrendResponse, error) {
	sales, err := u.saleRepo.FindByTenant(tenantID)
	if err != nil {
		return nil, err
	}

	trendMap := make(map[string]float64)
	now := time.Now()
	
	// Initialize last 7 days including today
	for i := 0; i < 7; i++ {
		dateStr := now.AddDate(0, 0, -i).Format("2006-01-02")
		trendMap[dateStr] = 0
	}

	for _, s := range sales {
		dateStr := s.CreatedAt.Format("2006-01-02")
		if _, ok := trendMap[dateStr]; ok {
			trendMap[dateStr] += s.GrandTotal
		}
	}

	result := make([]dto.SalesTrendResponse, 0)
	// Return in chronological order
	for i := 6; i >= 0; i-- {
		dateStr := now.AddDate(0, 0, -i).Format("2006-01-02")
		result = append(result, dto.SalesTrendResponse{
			Date:   dateStr,
			Amount: trendMap[dateStr],
		})
	}

	return result, nil
}
