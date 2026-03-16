package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/loan/dto"
	"github.com/koperasi-gresik/backend/internal/modules/loan/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type LoanHandler struct {
	service service.LoanService
}

func NewLoanHandler(service service.LoanService) *LoanHandler {
	return &LoanHandler{service: service}
}

func (h *LoanHandler) CreateProduct(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.LoanProductCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	product, err := h.service.CreateProduct(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, product, "Loan product created successfully")
}

func (h *LoanHandler) ListProducts(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	products, err := h.service.ListProducts(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, products)
}

func (h *LoanHandler) ApplyForLoan(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.LoanApplicationRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	loan, err := h.service.ApplyForLoan(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, loan, "Loan application submitted successfully")
}

func (h *LoanHandler) ApproveLoan(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	loanID, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID")
	}

	if err := h.service.ApproveLoan(c.Context(), orgID, uint(loanID)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "Loan approved successfully")
}

func (h *LoanHandler) RecordPayment(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	loanID, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID")
	}

	var req dto.LoanPaymentRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	payment, err := h.service.RecordPayment(c.Context(), orgID, uint(loanID), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, payment, "Loan payment recorded successfully")
}

func (h *LoanHandler) GetLoan(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	loanID, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID")
	}

	loan, err := h.service.GetLoanByID(c.Context(), orgID, uint(loanID))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, loan)
}

// RegisterRoutes registers the loan routes.
func RegisterRoutes(router fiber.Router, handler *LoanHandler, middlewares ...fiber.Handler) {
	group := router.Group("/loans", middlewares...)

	group.Post("/products", handler.CreateProduct)
	group.Get("/products", handler.ListProducts)

	group.Post("/apply", handler.ApplyForLoan)
	group.Get("/:id", handler.GetLoan)
	group.Post("/:id/approve", handler.ApproveLoan)
	group.Post("/:id/payments", handler.RecordPayment)
}
