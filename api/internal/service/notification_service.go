package service

import (
	"fmt"
	"lakoo/backend/internal/dto"
	"lakoo/backend/internal/repository"
)

type NotificationService interface {
	GetNotifications(tenantID string) ([]dto.NotificationResponse, error)
}

type notificationService struct {
	productRepo repository.ProductRepository
	saleRepo    repository.SaleRepository
}

func NewNotificationService(p repository.ProductRepository, s repository.SaleRepository) NotificationService {
	return &notificationService{
		productRepo: p,
		saleRepo:    s,
	}
}

func (u *notificationService) GetNotifications(tenantID string) ([]dto.NotificationResponse, error) {
	notifications := make([]dto.NotificationResponse, 0)

	// 1. Low Stock Alerts
	lowStock, err := u.productRepo.FindLowStock(tenantID)
	if err == nil {
		for _, p := range lowStock {
			notifications = append(notifications, dto.NotificationResponse{
				ID:        fmt.Sprintf("low-stock-%s", p.ID),
				Type:      "low_stock",
				Title:     "Stok Rendah",
				Message:   fmt.Sprintf("Produk %s tersisa %.0f %s", p.Name, p.StockQty, p.Unit),
				CreatedAt: p.UpdatedAt,
				IsRead:    false,
			})
		}
	}

	// 2. Recent Sales Alerts (e.g., last 5)
	recentSales, err := u.saleRepo.GetRecent(tenantID, 5)
	if err == nil {
		for _, s := range recentSales {
			notifications = append(notifications, dto.NotificationResponse{
				ID:        fmt.Sprintf("new-sale-%s", s.ID),
				Type:      "new_sale",
				Title:     "Penjualan Baru",
				Message:   fmt.Sprintf("Transaksi %s berhasil senilai Rp%.0f", s.InvoiceNo, s.GrandTotal),
				CreatedAt: s.CreatedAt,
				IsRead:    false,
			})
		}
	}

	return notifications, nil
}
