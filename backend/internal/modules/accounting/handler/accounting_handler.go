package handler

import (
	"strings"
	"time"

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
	accountType := c.Query("type")

	var accounts []dto.AccountResponse
	var err error

	if accountType != "" {
		// Use efficient database-level filtering
		accounts, err = h.service.ListAccountsByType(c.Context(), orgID, accountType)
		if err != nil {
			return response.InternalError(c, err.Error())
		}
	} else {
		accounts, err = h.service.ListAccounts(c.Context(), orgID)
		if err != nil {
			return response.InternalError(c, err.Error())
		}
	}

	return response.Success(c, accounts)
}

func (h *AccountingHandler) GetAccountByCode(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	code := c.Params("code")

	account, err := h.service.GetAccountByCode(c.Context(), orgID, code)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, account)
}

func (h *AccountingHandler) CreateJournalEntry(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.JournalEntryCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	// Check for idempotency key in header (trim whitespace)
	idempotencyKey := strings.TrimSpace(c.Get("Idempotency-Key"))

	var entry *dto.JournalEntryResponse
	var err error

	if idempotencyKey != "" {
		entry, err = h.service.CreateJournalEntryIdempotent(c.Context(), orgID, idempotencyKey, req)
	} else {
		entry, err = h.service.CreateJournalEntry(c.Context(), orgID, req)
	}

	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, entry, "Journal entry created successfully")
}

func (h *AccountingHandler) ListJournalEntries(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	// Check for filter parameters
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	sourceModule := c.Query("source_module")
	status := c.Query("status")
	referenceNum := c.Query("reference_number")

	// If any filter is present, use filtered list
	if startDate != "" || endDate != "" || sourceModule != "" || status != "" || referenceNum != "" {
		filter := dto.JournalEntryFilter{
			SourceModule: sourceModule,
			Status:       status,
			ReferenceNum: referenceNum,
		}

		// Parse dates if provided
		if startDate != "" {
			parsed, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				return response.BadRequest(c, "Invalid start_date format, expected YYYY-MM-DD")
			}
			filter.StartDate = &parsed
		}
		if endDate != "" {
			parsed, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				return response.BadRequest(c, "Invalid end_date format, expected YYYY-MM-DD")
			}
			filter.EndDate = &parsed
		}

		entries, err := h.service.ListJournalEntriesFiltered(c.Context(), orgID, filter)
		if err != nil {
			return response.InternalError(c, err.Error())
		}
		return response.Success(c, entries)
	}

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

func (h *AccountingHandler) ReverseJournalEntry(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid journal entry ID")
	}

	var req dto.ReverseJournalEntryRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	entry, err := h.service.ReverseJournalEntry(c.Context(), orgID, uint(id), req.Reason)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, entry, "Journal entry reversed successfully")
}

// RegisterRoutes registers the accounting routes.
func RegisterRoutes(router fiber.Router, handler *AccountingHandler, middlewares ...fiber.Handler) {
	group := router.Group("/accounting", middlewares...)

	// Chart of Accounts
	group.Post("/accounts", handler.CreateAccount)
	group.Get("/accounts", handler.ListAccounts)
	group.Get("/accounts/:code", handler.GetAccountByCode)

	// Journal Entries
	group.Post("/journal-entries", handler.CreateJournalEntry)
	group.Get("/journal-entries", handler.ListJournalEntries)
	group.Get("/journal-entries/:id", handler.GetJournalEntry)
	group.Post("/journal-entries/:id/reverse", handler.ReverseJournalEntry)
}
