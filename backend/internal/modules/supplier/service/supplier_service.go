package service

import (
	"context"
	"errors"

	"github.com/koperasi-gresik/backend/internal/modules/supplier/dto"
	"github.com/koperasi-gresik/backend/internal/modules/supplier/model"
	"github.com/koperasi-gresik/backend/internal/modules/supplier/repository"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
)

type SupplierService interface {
	Create(ctx context.Context, orgID uint, req dto.SupplierCreateRequest) (*dto.SupplierResponse, error)
	GetByID(ctx context.Context, orgID, id uint) (*dto.SupplierResponse, error)
	Update(ctx context.Context, orgID, id uint, req dto.SupplierUpdateRequest) (*dto.SupplierResponse, error)
	List(ctx context.Context, orgID uint, params pagination.Params) ([]dto.SupplierResponse, int64, error)
}

type supplierService struct {
	repo repository.SupplierRepository
}

func NewSupplierService(repo repository.SupplierRepository) SupplierService {
	return &supplierService{repo: repo}
}

func (s *supplierService) Create(ctx context.Context, orgID uint, req dto.SupplierCreateRequest) (*dto.SupplierResponse, error) {
	if _, err := s.repo.GetByCode(ctx, orgID, req.Code); err == nil {
		return nil, errors.New("supplier code already exists in this organization")
	}

	supplier := &model.Supplier{
		Code:        req.Code,
		Name:        req.Name,
		ContactName: req.ContactName,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		Status:      "active",
	}
	supplier.OrganizationID = orgID

	if err := s.repo.Create(ctx, supplier); err != nil {
		return nil, err
	}

	return s.mapToResponse(supplier), nil
}

func (s *supplierService) GetByID(ctx context.Context, orgID, id uint) (*dto.SupplierResponse, error) {
	supplier, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(supplier), nil
}

func (s *supplierService) Update(ctx context.Context, orgID, id uint, req dto.SupplierUpdateRequest) (*dto.SupplierResponse, error) {
	supplier, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}

	supplier.Name = req.Name
	supplier.ContactName = req.ContactName
	supplier.Phone = req.Phone
	supplier.Email = req.Email
	supplier.Address = req.Address

	if req.Status != "" {
		supplier.Status = req.Status
	}

	if err := s.repo.Update(ctx, supplier); err != nil {
		return nil, err
	}

	return s.mapToResponse(supplier), nil
}

func (s *supplierService) List(ctx context.Context, orgID uint, params pagination.Params) ([]dto.SupplierResponse, int64, error) {
	suppliers, total, err := s.repo.List(ctx, orgID, params)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.SupplierResponse
	for _, spl := range suppliers {
		res = append(res, *s.mapToResponse(&spl))
	}
	return res, total, nil
}

func (s *supplierService) mapToResponse(spl *model.Supplier) *dto.SupplierResponse {
	return &dto.SupplierResponse{
		ID:          spl.ID,
		Code:        spl.Code,
		Name:        spl.Name,
		ContactName: spl.ContactName,
		Phone:       spl.Phone,
		Email:       spl.Email,
		Address:     spl.Address,
		Status:      spl.Status,
		CreatedAt:   spl.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
