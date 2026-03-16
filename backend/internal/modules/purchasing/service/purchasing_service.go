package service

import (
	"context"
	"fmt"
	"time"

	inventorydto "github.com/koperasi-gresik/backend/internal/modules/inventory/dto"
	inventoryService "github.com/koperasi-gresik/backend/internal/modules/inventory/service"
	purchasingdto "github.com/koperasi-gresik/backend/internal/modules/purchasing/dto"
	"github.com/koperasi-gresik/backend/internal/modules/purchasing/model"
	"github.com/koperasi-gresik/backend/internal/modules/purchasing/repository"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type PurchasingService interface {
	CreatePO(ctx context.Context, orgID uint, req purchasingdto.PurchaseOrderCreateRequest) (*purchasingdto.PurchaseOrderResponse, error)
	GetPO(ctx context.Context, orgID, poID uint) (*purchasingdto.PurchaseOrderResponse, error)
	ReceivePO(ctx context.Context, orgID, poID uint) error
	ListPOs(ctx context.Context, orgID uint) ([]purchasingdto.PurchaseOrderResponse, error)
}

type purchasingService struct {
	repo             repository.PurchasingRepository
	inventoryService inventoryService.InventoryService
}

func NewPurchasingService(repo repository.PurchasingRepository, invService inventoryService.InventoryService) PurchasingService {
	return &purchasingService{
		repo:             repo,
		inventoryService: invService,
	}
}

func (s *purchasingService) CreatePO(ctx context.Context, orgID uint, req purchasingdto.PurchaseOrderCreateRequest) (*purchasingdto.PurchaseOrderResponse, error) {
	poNumber := utils.GenerateCode("PO")

	po := &model.PurchaseOrder{
		SupplierID: req.SupplierID,
		PONumber:   poNumber,
		Status:     "ordered",
		Notes:      req.Notes,
	}
	po.OrganizationID = orgID

	var totalAmount float64
	for _, itemReq := range req.Items {
		subtotal := itemReq.CostPrice * float64(itemReq.Quantity)
		item := model.PurchaseOrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			CostPrice: itemReq.CostPrice,
			Subtotal:  subtotal,
		}
		item.OrganizationID = orgID
		po.Items = append(po.Items, item)
		totalAmount += subtotal
	}

	po.TotalAmount = totalAmount
	po.Discount = req.Discount
	po.FinalAmount = totalAmount - req.Discount

	if err := s.repo.CreatePO(ctx, po); err != nil {
		return nil, err
	}

	return s.mapToResponse(po), nil
}

func (s *purchasingService) GetPO(ctx context.Context, orgID, poID uint) (*purchasingdto.PurchaseOrderResponse, error) {
	po, err := s.repo.GetPOByID(ctx, orgID, poID)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(po), nil
}

func (s *purchasingService) ReceivePO(ctx context.Context, orgID, poID uint) error {
	po, err := s.repo.GetPOByID(ctx, orgID, poID)
	if err != nil {
		return err
	}

	if po.Status == "received" {
		return fmt.Errorf("PO already received")
	}

	// 1. Update Inventory for each item
	for _, item := range po.Items {
		_, err := s.inventoryService.AdjustStock(ctx, orgID, item.ProductID, inventorydto.StockMovementRequest{
			Type:            "in",
			Quantity:        item.Quantity,
			Notes:           "PO Received: " + po.PONumber,
			RelatedEntity:   "purchase_order",
			RelatedEntityID: &po.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to update stock for product ID %d", item.ProductID)
		}
	}

	// 2. Update PO Status
	now := time.Now()
	po.Status = "received"
	po.ReceivedAt = &now

	return s.repo.UpdatePOStatus(ctx, po.ID, "received")
}

func (s *purchasingService) ListPOs(ctx context.Context, orgID uint) ([]purchasingdto.PurchaseOrderResponse, error) {
	pos, err := s.repo.ListPOs(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var res []purchasingdto.PurchaseOrderResponse
	for _, po := range pos {
		res = append(res, *s.mapToResponse(&po))
	}
	return res, nil
}

func (s *purchasingService) mapToResponse(po *model.PurchaseOrder) *purchasingdto.PurchaseOrderResponse {
	res := &purchasingdto.PurchaseOrderResponse{
		ID:            po.ID,
		PONumber:      po.PONumber,
		SupplierID:    po.SupplierID,
		TotalAmount:   po.TotalAmount,
		FinalAmount:   po.FinalAmount,
		PaymentStatus: po.PaymentStatus,
		Status:        po.Status,
		CreatedAt:     po.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	for _, item := range po.Items {
		res.Items = append(res.Items, purchasingdto.PurchaseOrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			CostPrice: item.CostPrice,
			Subtotal:  item.Subtotal,
		})
	}

	return res
}
