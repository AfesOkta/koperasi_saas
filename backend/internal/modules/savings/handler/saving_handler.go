package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/savings/dto"
	"github.com/koperasi-gresik/backend/internal/modules/savings/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type SavingHandler struct {
	service service.SavingService
}

func NewSavingHandler(service service.SavingService) *SavingHandler {
	return &SavingHandler{service: service}
}

func (h *SavingHandler) CreateProduct(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.SavingProductCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	product, err := h.service.CreateProduct(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, product, "Saving product created successfully")
}

func (h *SavingHandler) ListProducts(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	products, err := h.service.ListProducts(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, products)
}

func (h *SavingHandler) Deposit(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.SavingTransactionRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}
	req.Type = "deposit"

	txn, err := h.service.Deposit(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, txn, "Deposit successful")
}

func (h *SavingHandler) Withdraw(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.SavingTransactionRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}
	req.Type = "withdrawal"

	txn, err := h.service.Withdraw(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, txn, "Withdrawal successful")
}

func (h *SavingHandler) GetBalance(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	memberID, err := c.ParamsInt("memberId")
	if err != nil {
		return response.BadRequest(c, "Invalid member ID")
	}

	productID, err := c.ParamsInt("productId")
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	account, err := h.service.GetBalance(c.Context(), orgID, uint(memberID), uint(productID))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, account)
}

func (h *SavingHandler) GetHistory(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	accountID, err := c.ParamsInt("accountId")
	if err != nil {
		return response.BadRequest(c, "Invalid account ID")
	}

	history, err := h.service.GetTransactionHistory(c.Context(), orgID, uint(accountID))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, history)
}

// RegisterRoutes registers the saving routes.
func RegisterRoutes(router fiber.Router, handler *SavingHandler, middlewares ...fiber.Handler) {
	group := router.Group("/savings", middlewares...)

	// Products
	group.Post("/products", handler.CreateProduct)
	group.Get("/products", handler.ListProducts)

	// Transactions
	group.Post("/deposit", handler.Deposit)
	group.Post("/withdraw", handler.Withdraw)

	// Balance & History
	group.Get("/balance/:memberId/:productId", handler.GetBalance)
	group.Get("/history/:accountId", handler.GetHistory)
}
