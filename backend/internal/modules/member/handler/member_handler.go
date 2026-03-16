package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/member/dto"
	"github.com/koperasi-gresik/backend/internal/modules/member/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type MemberHandler struct {
	service service.MemberService
}

func NewMemberHandler(service service.MemberService) *MemberHandler {
	return &MemberHandler{service: service}
}

func (h *MemberHandler) Create(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.MemberCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	member, err := h.service.Create(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, member, "Member registered successfully")
}

func (h *MemberHandler) Get(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid member ID")
	}

	member, err := h.service.GetByID(c.Context(), orgID, uint(id))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, member)
}

func (h *MemberHandler) List(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	params := pagination.Parse(c)

	members, total, err := h.service.List(c.Context(), orgID, params)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Paginated(c, members, response.Meta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: params.TotalPages(total),
	})
}

func (h *MemberHandler) UpdateStatus(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid member ID")
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=active inactive"`
	}
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	if err := h.service.UpdateStatus(c.Context(), orgID, uint(id), req.Status); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "Member status updated successfully")
}

func (h *MemberHandler) UploadDocument(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid member ID")
	}

	docType := c.FormValue("type")
	if docType == "" {
		return response.BadRequest(c, "Document type is required")
	}

	fileHeader, err := c.FormFile("document")
	if err != nil {
		return response.BadRequest(c, "Document file is required")
	}

	// Pseudo logic: Save to disk or S3
	// For MVP: simply save locally to a ./uploads dir
	fileName := fmt.Sprintf("org_%d_mem_%d_%s_%s", orgID, id, docType, fileHeader.Filename)
	savePath := fmt.Sprintf("./uploads/%s", fileName)

	if err := c.SaveFile(fileHeader, savePath); err != nil {
		return response.InternalError(c, "Failed to ave uploaded file")
	}

	fileURL := fmt.Sprintf("/uploads/%s", fileName) // In reality this would be S3 URL

	if err := h.service.UploadDocument(c.Context(), orgID, uint(id), docType, fileURL); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, fiber.Map{"file_url": fileURL}, "Document uploaded successfully")
}

// RegisterRoutes registers the member routes.
func RegisterRoutes(router fiber.Router, handler *MemberHandler, middlewares ...fiber.Handler) {
	group := router.Group("/members", middlewares...)
	group.Post("/", handler.Create)
	group.Get("/", handler.List)
	group.Get("/:id", handler.Get)
	group.Patch("/:id/status", handler.UpdateStatus)
	group.Post("/:id/documents", handler.UploadDocument)
}
