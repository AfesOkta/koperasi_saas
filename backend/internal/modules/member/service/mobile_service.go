package service

import (
	"context"
	"fmt"

	"github.com/koperasi-gresik/backend/internal/modules/member/dto"
	"github.com/koperasi-gresik/backend/internal/modules/member/repository"
	savingRepo "github.com/koperasi-gresik/backend/internal/modules/savings/repository"
	loanRepo "github.com/koperasi-gresik/backend/internal/modules/loan/repository"
)

type MobileService interface {
	GetMemberDashboard(ctx context.Context, orgID, userID uint) (*dto.MemberDashboardResponse, error)
}

type mobileService struct {
	memberRepo repository.MemberRepository
	savingRepo savingRepo.SavingRepository
	loanRepo   loanRepo.LoanRepository
}

func NewMobileService(
	memberRepo repository.MemberRepository,
	savingRepo savingRepo.SavingRepository,
	loanRepo loanRepo.LoanRepository,
) MobileService {
	return &mobileService{
		memberRepo: memberRepo,
		savingRepo: savingRepo,
		loanRepo:   loanRepo,
	}
}

func (s *mobileService) GetMemberDashboard(ctx context.Context, orgID, userID uint) (*dto.MemberDashboardResponse, error) {
	// 1. Get Member by UserID
	member, err := s.memberRepo.GetByUserID(ctx, orgID, userID)
	if err != nil {
		return nil, fmt.Errorf("member profile not found")
	}

	// 2. Fetch Total Savings
	savings, _ := s.savingRepo.ListAccountsByMember(ctx, orgID, member.ID)
	var totalSavings float64
	for _, acc := range savings {
		totalSavings += acc.Balance
	}

	// 3. Fetch Total Outstanding Loans
	loans, _ := s.loanRepo.ListLoansByMember(ctx, orgID, member.ID)
	var totalLoans float64
	for _, l := range loans {
		if l.Status == "active" {
			totalLoans += l.Outstanding
		}
	}

	// 4. Generate QR Code Data (Basic format: orgID:memberID:NIK)
	qrData := fmt.Sprintf("MEMBER:%d:%d:%s", orgID, member.ID, member.NIK)

	return &dto.MemberDashboardResponse{
		Profile: dto.MemberResponse{
			ID:           member.ID,
			Name:         member.Name,
			MemberNumber: member.MemberNumber,
			Status:       member.Status,
		},
		TotalSavings: totalSavings,
		TotalLoans:    totalLoans,
		QRCode:       qrData,
	}, nil
}
