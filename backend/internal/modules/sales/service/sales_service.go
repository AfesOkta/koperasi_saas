package service

import (
	"context"
	"fmt"

	inventorydto "github.com/koperasi-gresik/backend/internal/modules/inventory/dto"
	inventoryService "github.com/koperasi-gresik/backend/internal/modules/inventory/service"
	salesdto "github.com/koperasi-gresik/backend/internal/modules/sales/dto"
	"github.com/koperasi-gresik/backend/internal/modules/sales/model"
	"github.com/koperasi-gresik/backend/internal/modules/sales/repository"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type SalesService interface {
	CreateOrder(ctx context.Context, orgID, cashierID uint, req salesdto.OrderCreateRequest) (*salesdto.OrderResponse, error)
	GetOrder(ctx context.Context, orgID, orderID uint) (*salesdto.OrderResponse, error)
	ListOrders(ctx context.Context, orgID uint) ([]salesdto.OrderResponse, error)
}

type salesService struct {
	repo             repository.SalesRepository
	inventoryService inventoryService.InventoryService
}

func NewSalesService(repo repository.SalesRepository, invService inventoryService.InventoryService) SalesService {
	return &salesService{
		repo:             repo,
		inventoryService: invService,
	}
}

func (s *salesService) CreateOrder(ctx context.Context, orgID, cashierID uint, req salesdto.OrderCreateRequest) (*salesdto.OrderResponse, error) {
	orderID := utils.GenerateCode("ORD")

	order := &model.Order{
		OrderID:   orderID,
		MemberID:  req.MemberID,
		CashierID: cashierID,
		Status:    "completed",
	}
	order.OrganizationID = orgID

	var totalAmount float64
	var orderItems []model.OrderItem

	// 1. Process Items and validate inventory
	for _, itemReq := range req.Items {
		// Get Product Details
		product, err := s.inventoryService.GetProduct(ctx, orgID, itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found", itemReq.ProductID)
		}

		if product.Stock < itemReq.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s (%s)", product.Name, product.SKU)
		}

		subtotal := product.Price * float64(itemReq.Quantity)
		lineTotal := subtotal - itemReq.Discount

		orderItem := model.OrderItem{
			ProductID:   itemReq.ProductID,
			Quantity:    itemReq.Quantity,
			UnitPrice:   product.Price,
			Subtotal:    subtotal,
			Discount:    itemReq.Discount,
			TotalAmount: lineTotal,
		}
		orderItem.OrganizationID = orgID
		orderItems = append(orderItems, orderItem)

		totalAmount += lineTotal

		// 2. Reduce Stock
		_, err = s.inventoryService.AdjustStock(ctx, orgID, itemReq.ProductID, inventorydto.StockMovementRequest{
			Type:          "out",
			Quantity:      itemReq.Quantity,
			Notes:         "Sales Order " + orderID,
			RelatedEntity: "sales_order",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to reduce stock for product %s", product.Name)
		}
	}

	order.TotalAmount = totalAmount
	order.Discount = req.Discount
	order.FinalAmount = totalAmount - req.Discount
	order.Items = orderItems

	// 3. Process Payments
	var paidAmount float64
	for _, payReq := range req.Payments {
		payment := model.OrderPayment{
			PaymentMethod:  payReq.Method,
			Amount:         payReq.Amount,
			ReferenceToken: payReq.ReferenceToken,
			Status:         "completed",
		}
		payment.OrganizationID = orgID
		order.Payments = append(order.Payments, payment)
		paidAmount += payReq.Amount
	}

	if paidAmount >= order.FinalAmount {
		order.PaymentStatus = "paid"
	} else if paidAmount > 0 {
		order.PaymentStatus = "partial"
	} else {
		order.PaymentStatus = "unpaid"
	}

	// 4. Save Order
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	// TODO: Integrate with Accounting & Cash modules to record inflows

	return s.mapOrderToResponse(order), nil
}

func (s *salesService) GetOrder(ctx context.Context, orgID, orderID uint) (*salesdto.OrderResponse, error) {
	order, err := s.repo.GetOrderByID(ctx, orgID, orderID)
	if err != nil {
		return nil, err
	}
	return s.mapOrderToResponse(order), nil
}

func (s *salesService) ListOrders(ctx context.Context, orgID uint) ([]salesdto.OrderResponse, error) {
	orders, err := s.repo.ListOrders(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var res []salesdto.OrderResponse
	for _, o := range orders {
		res = append(res, *s.mapOrderToResponse(&o))
	}
	return res, nil
}

// Mappers
func (s *salesService) mapOrderToResponse(o *model.Order) *salesdto.OrderResponse {
	res := &salesdto.OrderResponse{
		ID:            o.ID,
		OrderID:       o.OrderID,
		MemberID:      o.MemberID,
		TotalAmount:   o.TotalAmount,
		Discount:      o.Discount,
		TaxAmount:     o.TaxAmount,
		FinalAmount:   o.FinalAmount,
		PaymentStatus: o.PaymentStatus,
		Status:        o.Status,
		CashierID:     o.CashierID,
		CreatedAt:     o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	for _, item := range o.Items {
		res.Items = append(res.Items, salesdto.OrderItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.Subtotal,
			Discount:    item.Discount,
			TotalAmount: item.TotalAmount,
		})
	}

	for _, p := range o.Payments {
		res.Payments = append(res.Payments, salesdto.OrderPaymentResponse{
			ID:            p.ID,
			PaymentMethod: p.PaymentMethod,
			Amount:        p.Amount,
			Status:        p.Status,
		})
	}

	return res
}
