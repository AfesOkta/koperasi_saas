package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/member/model"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

type MemberRepository interface {
	Create(ctx context.Context, member *model.Member) error
	GetByID(ctx context.Context, orgID, id uint) (*model.Member, error)
	GetByNumber(ctx context.Context, orgID uint, number string) (*model.Member, error)
	GetByNIK(ctx context.Context, orgID uint, nik string) (*model.Member, error)
	Update(ctx context.Context, member *model.Member) error
	List(ctx context.Context, orgID uint, params pagination.Params) ([]model.Member, int64, error)

	AddDocument(ctx context.Context, doc *model.MemberDocument) error
	AddCard(ctx context.Context, card *model.MemberCard) error
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) Create(ctx context.Context, member *model.Member) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *memberRepository) GetByID(ctx context.Context, orgID, id uint) (*model.Member, error) {
	var member model.Member
	err := r.db.WithContext(ctx).
		Preload("Documents").
		Preload("Cards").
		Where("organization_id = ?", orgID).
		First(&member, id).Error
	return &member, err
}

func (r *memberRepository) GetByNumber(ctx context.Context, orgID uint, number string) (*model.Member, error) {
	var member model.Member
	err := r.db.WithContext(ctx).
		Preload("Cards").
		Where("organization_id = ? AND member_number = ?", orgID, number).
		First(&member).Error
	return &member, err
}

func (r *memberRepository) GetByNIK(ctx context.Context, orgID uint, nik string) (*model.Member, error) {
	var member model.Member
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND nik = ?", orgID, nik).
		First(&member).Error
	return &member, err
}

func (r *memberRepository) Update(ctx context.Context, member *model.Member) error {
	return r.db.WithContext(ctx).Save(member).Error
}

func (r *memberRepository) List(ctx context.Context, orgID uint, params pagination.Params) ([]model.Member, int64, error) {
	var members []model.Member
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Member{}).Where("organization_id = ?", orgID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(params.Scope()).Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

func (r *memberRepository) AddDocument(ctx context.Context, doc *model.MemberDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *memberRepository) AddCard(ctx context.Context, card *model.MemberCard) error {
	return r.db.WithContext(ctx).Create(card).Error
}
