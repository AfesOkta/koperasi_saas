package service

import (
	"context"
	"errors"
	"fmt"

	iamDTO "github.com/koperasi-gresik/backend/internal/modules/iam/dto"
	iamService "github.com/koperasi-gresik/backend/internal/modules/iam/service"
	"github.com/koperasi-gresik/backend/internal/modules/member/dto"
	"github.com/koperasi-gresik/backend/internal/modules/member/model"
	"github.com/koperasi-gresik/backend/internal/modules/member/repository"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type MemberService interface {
	Create(ctx context.Context, orgID uint, req dto.MemberCreateRequest) (*dto.MemberResponse, error)
	GetByID(ctx context.Context, orgID, id uint) (*dto.MemberResponse, error)
	List(ctx context.Context, orgID uint, params pagination.Params) ([]dto.MemberResponse, int64, error)
	UpdateStatus(ctx context.Context, orgID, id uint, status string) error

	UploadDocument(ctx context.Context, orgID, memberID uint, docType, fileURL string) error
}

type memberService struct {
	repo        repository.MemberRepository
	authService iamService.AuthService // For creating optional IAM logins
}

func NewMemberService(repo repository.MemberRepository, authService iamService.AuthService) MemberService {
	return &memberService{
		repo:        repo,
		authService: authService,
	}
}

func (s *memberService) Create(ctx context.Context, orgID uint, req dto.MemberCreateRequest) (*dto.MemberResponse, error) {
	// Check NIK uniqueness within org
	if _, err := s.repo.GetByNIK(ctx, orgID, req.NIK); err == nil {
		return nil, errors.New("member with this NIK already exists")
	}

	member := &model.Member{
		MemberNumber: utils.GenerateCode("MBR"),
		Name:         req.Name,
		NIK:          req.NIK,
		Address:      req.Address,
		Phone:        req.Phone,
		Status:       "pending",
	}
	member.OrganizationID = orgID

	// Create IAM Login if requested
	if req.CreateSystem && req.Email != "" {
		userReq := iamDTO.UserCreateRequest{
			Name:     req.Name,
			Email:    req.Email,
			Password: utils.GenerateCode("PASS"), // Generate temporary password
			Phone:    req.Phone,
			RoleID:   0, // Member role ID logic goes here eventually
		}

		userRes, err := s.authService.Register(ctx, userReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create system user: %w", err)
		}
		member.UserID = userRes.ID
	}

	if err := s.repo.Create(ctx, member); err != nil {
		return nil, errors.New("failed to register member")
	}

	return s.mapToResponse(member), nil
}

func (s *memberService) GetByID(ctx context.Context, orgID, id uint) (*dto.MemberResponse, error) {
	member, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, errors.New("member not found")
	}

	return s.mapToResponse(member), nil
}

func (s *memberService) List(ctx context.Context, orgID uint, params pagination.Params) ([]dto.MemberResponse, int64, error) {
	members, total, err := s.repo.List(ctx, orgID, params)
	if err != nil {
		return nil, 0, errors.New("failed to list members")
	}

	responses := make([]dto.MemberResponse, len(members))
	for i, m := range members {
		responses[i] = *s.mapToResponse(&m)
	}

	return responses, total, nil
}

func (s *memberService) UpdateStatus(ctx context.Context, orgID, id uint, status string) error {
	member, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return errors.New("member not found")
	}

	member.Status = status
	if err := s.repo.Update(ctx, member); err != nil {
		return errors.New("failed to update member status")
	}

	return nil
}

func (s *memberService) UploadDocument(ctx context.Context, orgID, memberID uint, docType, fileURL string) error {
	// Verify member belongs to org
	if _, err := s.repo.GetByID(ctx, orgID, memberID); err != nil {
		return errors.New("member not found")
	}

	doc := &model.MemberDocument{
		MemberID: memberID,
		Type:     docType,
		FileURL:  fileURL,
	}
	doc.OrganizationID = orgID

	if err := s.repo.AddDocument(ctx, doc); err != nil {
		return errors.New("failed to save document record")
	}

	return nil
}

func (s *memberService) mapToResponse(m *model.Member) *dto.MemberResponse {
	resp := &dto.MemberResponse{
		ID:           m.ID,
		UserID:       m.UserID,
		MemberNumber: m.MemberNumber,
		Name:         m.Name,
		NIK:          m.NIK,
		Address:      m.Address,
		Phone:        m.Phone,
		Status:       m.Status,
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if len(m.Documents) > 0 {
		for _, d := range m.Documents {
			resp.Documents = append(resp.Documents, dto.MemberDocumentResponse{
				ID:      d.ID,
				Type:    d.Type,
				FileURL: d.FileURL,
			})
		}
	}

	if len(m.Cards) > 0 {
		for _, c := range m.Cards {
			resp.Cards = append(resp.Cards, dto.MemberCardResponse{
				ID:         c.ID,
				CardNumber: c.CardNumber,
				Status:     c.Status,
			})
		}
	}

	return resp
}
