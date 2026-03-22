package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	inventorydto "github.com/koperasi-gresik/backend/internal/modules/inventory/dto"
	inventoryService "github.com/koperasi-gresik/backend/internal/modules/inventory/service"
	"github.com/koperasi-gresik/backend/internal/modules/pos/dto"
	"github.com/koperasi-gresik/backend/internal/modules/pos/model"
	"github.com/koperasi-gresik/backend/internal/modules/pos/repository"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type POSService interface {
	// Shift Management
	OpenShift(ctx context.Context, orgID, cashierID uint, startBalance float64) (*dto.ShiftResponse, error)
	GetActiveShift(ctx context.Context, cashierID uint) (*dto.ShiftResponse, error)

	// Order Management
	CreateOrder(ctx context.Context, orgID uint, req dto.OrderCreateRequest) (*dto.OrderResponse, error)
	GetOrder(ctx context.Context, orgID, orderID uint) (*dto.OrderResponse, error)

	// KDS Management
	UpdateKDSStatus(ctx context.Context, orgID, itemID uint, status string) error
	GetKDSItems(ctx context.Context, orgID uint, statuses []string) ([]dto.KDSItemResponse, error)

	// Receipt Management
	GenerateReceipt(ctx context.Context, orgID, orderID uint) (string, error)
}

type posService struct {
	repo             repository.POSRepository
	inventoryService inventoryService.InventoryService
}

func NewPOSService(repo repository.POSRepository, invService inventoryService.InventoryService) POSService {
	return &posService{
		repo:             repo,
		inventoryService: invService,
	}
}

func (s *posService) OpenShift(ctx context.Context, orgID, cashierID uint, startBalance float64) (*dto.ShiftResponse, error) {
	shift := &model.POSShift{
		CashierID:    cashierID,
		StartTime:    time.Now(),
		StartBalance: startBalance,
		Status:       "open",
	}
	shift.OrganizationID = orgID

	if err := s.repo.OpenShift(ctx, shift); err != nil {
		return nil, err
	}

	return &dto.ShiftResponse{
		ID:           shift.ID,
		CashierID:    shift.CashierID,
		StartTime:    shift.StartTime.Format(time.RFC3339),
		Status:       shift.Status,
		StartBalance: shift.StartBalance,
	}, nil
}

func (s *posService) GetActiveShift(ctx context.Context, cashierID uint) (*dto.ShiftResponse, error) {
	shift, err := s.repo.GetActiveShift(ctx, cashierID)
	if err != nil {
		return nil, err
	}
	return &dto.ShiftResponse{
		ID:           shift.ID,
		CashierID:    shift.CashierID,
		StartTime:    shift.StartTime.Format(time.RFC3339),
		Status:       shift.Status,
		StartBalance: shift.StartBalance,
	}, nil
}

func (s *posService) CreateOrder(ctx context.Context, orgID uint, req dto.OrderCreateRequest) (*dto.OrderResponse, error) {
	refNumber := utils.GenerateCode("POS")
	
	order := &model.POSOrder{
		ShiftID:         req.ShiftID,
		ReferenceNumber: refNumber,
		PaymentMethod:   req.PaymentMethod,
		Status:          "completed",
		Notes:           req.Notes,
	}
	order.OrganizationID = orgID

	var totalAmount float64
	for _, itemReq := range req.Items {
		// 1. Get Product Data from Inventory module
		product, err := s.inventoryService.GetProduct(ctx, orgID, itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to find product %d: %w", itemReq.ProductID, err)
		}

		subtotal := product.Price * float64(itemReq.Quantity)
		
		item := model.POSOrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: product.Price,
			Subtotal:  subtotal,
			KDSStatus: "pending",
			Notes:     itemReq.Notes,
		}
		item.OrganizationID = orgID
		order.Items = append(order.Items, item)
		totalAmount += subtotal

		// 2. Deduct Stock (Main Warehouse for now)
		_, err = s.inventoryService.AdjustStock(ctx, orgID, itemReq.ProductID, inventorydto.StockMovementRequest{
			Type:            "out",
			Quantity:        itemReq.Quantity,
			Notes:           "POS Order: " + refNumber,
			RelatedEntity:   "pos_order",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to deduct stock for product %d: %w", itemReq.ProductID, err)
		}
	}

	order.TotalAmount = totalAmount
	order.FinalAmount = totalAmount // Simplified for MVP (no tax/discount in req yet)

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return s.mapOrderToResponse(order), nil
}

func (s *posService) GetOrder(ctx context.Context, orgID, orderID uint) (*dto.OrderResponse, error) {
	order, err := s.repo.GetOrderByID(ctx, orgID, orderID)
	if err != nil {
		return nil, err
	}
	return s.mapOrderToResponse(order), nil
}

func (s *posService) UpdateKDSStatus(ctx context.Context, orgID, itemID uint, status string) error {
	return s.repo.UpdateOrderItemKDSStatus(ctx, orgID, itemID, status)
}

func (s *posService) GetKDSItems(ctx context.Context, orgID uint, statuses []string) ([]dto.KDSItemResponse, error) {
	items, err := s.repo.GetKDSItems(ctx, orgID, statuses)
	if err != nil {
		return nil, err
	}

	var res []dto.KDSItemResponse
	for _, item := range items {
		res = append(res, dto.KDSItemResponse{
			ID:             item.ID,
			OrderID:        item.OrderID,
			OrderReference: item.Order.ReferenceNumber,
			ProductName:    "Product #" + fmt.Sprint(item.ProductID), // Need product name preload or cache
			Quantity:       item.Quantity,
			KDSStatus:      item.KDSStatus,
			Notes:          item.Notes,
			CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		})
	}
	return res, nil
}

func (s *posService) GenerateReceipt(ctx context.Context, orgID, orderID uint) (string, error) {
	order, err := s.repo.GetOrderByID(ctx, orgID, orderID)
	if err != nil {
		return "", err
	}

	builder := utils.NewReceiptBuilder()
	builder.AddCentered("KOPERASI G-SAAS")
	builder.AddCentered("Ref: " + order.ReferenceNumber)
	builder.AddDivider()

	for _, item := range order.Items {
		name := fmt.Sprintf("%dx Product #%d", item.Quantity, item.ProductID)
		price := utils.FormatPrice(item.Subtotal)
		builder.AddKeyValue(name, price)
	}

	builder.AddDivider()
	builder.AddKeyValue("TOTAL", utils.FormatPrice(order.FinalAmount))
	builder.AddCentered("Payment: " + strings.ToUpper(order.PaymentMethod))
	builder.AddDivider()
	builder.AddCentered("Please visit again!")
	builder.AddCentered(time.Now().Format("2006-01-02 15:04"))

	return builder.Cut(), nil
}

func (s *posService) mapOrderToResponse(order *model.POSOrder) *dto.OrderResponse {
	res := &dto.OrderResponse{
		ID:              order.ID,
		ReferenceNumber: order.ReferenceNumber,
		TotalAmount:     order.TotalAmount,
		FinalAmount:     order.FinalAmount,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt.Format(time.RFC3339),
	}

	for _, item := range order.Items {
		res.Items = append(res.Items, dto.OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
			KDSStatus: item.KDSStatus,
			Notes:     item.Notes,
		})
	}
	return res
}
