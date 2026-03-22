package service

import (
	"context"
	"fmt"
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/inventory/dto"
	"github.com/koperasi-gresik/backend/internal/modules/inventory/model"
	"github.com/koperasi-gresik/backend/internal/modules/inventory/repository"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type WarehouseService interface {
	CreateWarehouse(ctx context.Context, orgID uint, req dto.WarehouseCreateRequest) (*dto.WarehouseResponse, error)
	ListWarehouses(ctx context.Context, orgID uint) ([]dto.WarehouseResponse, error)
	
	// Stock Transfers
	InitiateTransfer(ctx context.Context, orgID uint, req dto.TransferCreateRequest) (*dto.TransferResponse, error)
	ShipTransfer(ctx context.Context, orgID, transferID uint) error
	ReceiveTransfer(ctx context.Context, orgID, transferID uint) error
}

type warehouseService struct {
	repo          repository.WarehouseRepository
	inventoryRepo repository.InventoryRepository
}

func NewWarehouseService(repo repository.WarehouseRepository, inventoryRepo repository.InventoryRepository) WarehouseService {
	return &warehouseService{
		repo:          repo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *warehouseService) CreateWarehouse(ctx context.Context, orgID uint, req dto.WarehouseCreateRequest) (*dto.WarehouseResponse, error) {
	warehouse := &model.Warehouse{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		IsActive:    true,
	}
	warehouse.OrganizationID = orgID

	if err := s.repo.CreateWarehouse(ctx, warehouse); err != nil {
		return nil, err
	}

	return s.mapToResponse(warehouse), nil
}

func (s *warehouseService) ListWarehouses(ctx context.Context, orgID uint) ([]dto.WarehouseResponse, error) {
	warehouses, err := s.repo.ListWarehouses(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var res []dto.WarehouseResponse
	for _, w := range warehouses {
		res = append(res, *s.mapToResponse(&w))
	}
	return res, nil
}

func (s *warehouseService) InitiateTransfer(ctx context.Context, orgID uint, req dto.TransferCreateRequest) (*dto.TransferResponse, error) {
	// 1. Check Source Stock Availability
	item, err := s.repo.GetWarehouseItem(ctx, orgID, req.FromWarehouseID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("source warehouse does not have this product: %w", err)
	}

	if item.Quantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock in source warehouse (available: %d)", item.Quantity)
	}

	// 2. Create Transfer in 'pending' status
	transfer := &model.StockTransfer{
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
		ReferenceNumber: utils.GenerateCode("TRF"),
		FromWarehouseID: req.FromWarehouseID,
		ToWarehouseID:   req.ToWarehouseID,
		Status:          "pending",
		Notes:           req.Notes,
	}
	transfer.OrganizationID = orgID

	if err := s.repo.CreateTransfer(ctx, transfer); err != nil {
		return nil, err
	}

	return s.mapTransferToResponse(transfer), nil
}

func (s *warehouseService) ShipTransfer(ctx context.Context, orgID, transferID uint) error {
	transfer, err := s.repo.GetTransferByID(ctx, orgID, transferID)
	if err != nil {
		return err
	}

	if transfer.Status != "pending" {
		return fmt.Errorf("transfer cannot be shipped (current status: %s)", transfer.Status)
	}

	// 1. Move stock from Source Warehouse to 'Transit'
	// This subtracts from Source Warehouse balance
	now := time.Now()
	movement := &model.StockMovement{
		ProductID:       transfer.ProductID,
		WarehouseID:     transfer.FromWarehouseID,
		ReferenceNumber: transfer.ReferenceNumber,
		Type:            "transfer_out",
		Quantity:        transfer.Quantity,
		Notes:           "Shipping transfer to " + transfer.ToWarehouse.Name,
		RelatedEntity:   "stock_transfer",
		RelatedEntityID: &transfer.ID,
	}
	movement.OrganizationID = orgID

	if err := s.inventoryRepo.AdjustStock(ctx, transfer.FromWarehouseID, transfer.ProductID, movement); err != nil {
		return err
	}

	// 2. Update Transfer Status
	transfer.Status = "shipped"
	transfer.ShippedAt = &now
	return s.repo.UpdateTransfer(ctx, transfer)
}

func (s *warehouseService) ReceiveTransfer(ctx context.Context, orgID, transferID uint) error {
	transfer, err := s.repo.GetTransferByID(ctx, orgID, transferID)
	if err != nil {
		return err
	}

	if transfer.Status != "shipped" {
		return fmt.Errorf("transfer cannot be received (current status: %s)", transfer.Status)
	}

	// 1. Add stock to Destination Warehouse
	now := time.Now()
	movement := &model.StockMovement{
		ProductID:       transfer.ProductID,
		WarehouseID:     transfer.ToWarehouseID,
		ReferenceNumber: transfer.ReferenceNumber,
		Type:            "transfer_in",
		Quantity:        transfer.Quantity,
		Notes:           "Received transfer from " + transfer.FromWarehouse.Name,
		RelatedEntity:   "stock_transfer",
		RelatedEntityID: &transfer.ID,
	}
	movement.OrganizationID = orgID

	if err := s.inventoryRepo.AdjustStock(ctx, transfer.ToWarehouseID, transfer.ProductID, movement); err != nil {
		return err
	}

	// 2. Update Transfer Status
	transfer.Status = "received"
	transfer.ReceivedAt = &now
	return s.repo.UpdateTransfer(ctx, transfer)
}

func (s *warehouseService) mapToResponse(w *model.Warehouse) *dto.WarehouseResponse {
	return &dto.WarehouseResponse{
		ID:          w.ID,
		Code:        w.Code,
		Name:        w.Name,
		Description: w.Description,
		Address:     w.Address,
		IsActive:    w.IsActive,
	}
}

func (s *warehouseService) mapTransferToResponse(t *model.StockTransfer) *dto.TransferResponse {
	return &dto.TransferResponse{
		ID:              t.ID,
		ReferenceNumber: t.ReferenceNumber,
		FromWarehouse:   *s.mapToResponse(&t.FromWarehouse),
		ToWarehouse:     *s.mapToResponse(&t.ToWarehouse),
		Status:          t.Status,
		Notes:           t.Notes,
		CreatedAt:       t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
