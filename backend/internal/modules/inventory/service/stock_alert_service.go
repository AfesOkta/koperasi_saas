package service

import (
	"context"
	"fmt"
	"log"

	"github.com/koperasi-gresik/backend/internal/modules/inventory/repository"
	"github.com/koperasi-gresik/backend/internal/shared/event"
)

type StockAlertService interface {
	CheckLowStock(ctx context.Context, orgID uint) error
}

type stockAlertService struct {
	repo      repository.WarehouseRepository
	publisher event.Publisher
}

func NewStockAlertService(repo repository.WarehouseRepository, pub event.Publisher) StockAlertService {
	return &stockAlertService{
		repo:      repo,
		publisher: pub,
	}
}

func (s *stockAlertService) CheckLowStock(ctx context.Context, orgID uint) error {
	// 1. Get all low stock items for this organization
	items, err := s.repo.GetLowStockItems(ctx, orgID)
	if err != nil {
		return fmt.Errorf("failed to get low stock items: %w", err)
	}

	for _, item := range items {
		log.Printf("⚠️  [LowStock] Warehouse: %s, Product: %d, Qty: %d, Min: %d", 
			item.Warehouse.Name, item.ProductID, item.Quantity, item.MinStock)
		
		// 2. Emit Event
		evt := event.Event{
			Type:           event.EventStockLow,
			OrganizationID: orgID,
			AggregateID:    item.ProductID,
			Timestamp:      0, // Will be set by publisher/bus if needed
			Payload: map[string]interface{}{
				"warehouse_id":   item.WarehouseID,
				"warehouse_name": item.Warehouse.Name,
				"quantity":       item.Quantity,
				"min_stock":      item.MinStock,
			},
		}
		
		if err := s.publisher.Publish(ctx, evt); err != nil {
			log.Printf("❌ Failed to publish low stock event: %v", err)
		}
	}

	return nil
}
