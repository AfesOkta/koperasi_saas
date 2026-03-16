package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/accounting/dto"
	"github.com/koperasi-gresik/backend/internal/modules/accounting/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type AccountingHandler struct {
	service service.AccountingService
}

func NewAccountingHandler(service service.AccountingService) *AccountingHandler {
	return &AccountingHandler{service: service}
}

func (h *AccountingHandler) CreateAccount(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.AccountCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	account, err := h.service.CreateAccount(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, account, "Chart of Account created successfully")
}

func (h *AccountingHandler) ListAccounts(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	accounts, err := h.service.ListAccounts(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, accounts)
}

func (h *AccountingHandler) CreateJournalEntry(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.JournalEntryCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	entry, err := h.service.CreateJournalEntry(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, entry, "Journal entry created successfully")
}

func (h *AccountingHandler) ListJournalEntries(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	entries, err := h.service.ListJournalEntries(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, entries)
}

func (h *AccountingHandler) GetJournalEntry(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid journal entry ID")
	}

	entry, err := h.service.GetJournalEntryByID(c.Context(), orgID, uint(id))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, entry)
}

// RegisterRoutes registers the accounting routes.
func RegisterRoutes(router fiber.Router, handler *AccountingHandler, middlewares ...fiber.Handler) {
	group := router.Group("/accounting", middlewares...)

	// Chart of Accounts
	group.Post("/accounts", handler.CreateAccount)
	group.Get("/accounts", handler.ListAccounts)

	// Journal Entries
	group.Post("/journal-entries", handler.CreateJournalEntry)
	group.Get("/journal-entries", handler.ListJournalEntries)
	group.Get("/journal-entries/:id", handler.GetJournalEntry)
}
