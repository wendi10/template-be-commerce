package mysql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/pkg/apperrors"
)

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) *promoRepository {
	return &promoRepository{db: db}
}

func (r *promoRepository) Create(ctx context.Context, promo *domain.PromoCode) error {
	if promo.ID == uuid.Nil {
		promo.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(promo).Error; err != nil {
		if isUniqueViolation(err) {
			return apperrors.ErrInvalidPromo // promo code already exists
		}
		return err
	}
	return nil
}

func (r *promoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PromoCode, error) {
	var p domain.PromoCode
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *promoRepository) GetByCode(ctx context.Context, code string) (*domain.PromoCode, error) {
	var p domain.PromoCode
	err := r.db.WithContext(ctx).
		Where("code = ? AND deleted_at IS NULL", code).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *promoRepository) Update(ctx context.Context, promo *domain.PromoCode) error {
	return r.db.WithContext(ctx).
		Omit("used_count"). // never allow direct update of usage count
		Save(promo).Error
}

func (r *promoRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.PromoCode{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *promoRepository) IncrementUsage(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&domain.PromoCode{}).
		Where("id = ? AND (usage_limit = 0 OR used_count < usage_limit) AND deleted_at IS NULL", id).
		Update("used_count", gorm.Expr("used_count + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrPromoUsageLimitReached
	}
	return nil
}

func (r *promoRepository) List(ctx context.Context) ([]domain.PromoCode, error) {
	var promos []domain.PromoCode
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&promos).Error
	return promos, err
}
