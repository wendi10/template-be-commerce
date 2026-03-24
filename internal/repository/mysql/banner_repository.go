package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/pkg/apperrors"
)

type bannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) *bannerRepository {
	return &bannerRepository{db: db}
}

func (r *bannerRepository) Create(ctx context.Context, banner *domain.Banner) error {
	if banner.ID == uuid.Nil {
		banner.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(banner).Error
}

func (r *bannerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Banner, error) {
	var b domain.Banner
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *bannerRepository) Update(ctx context.Context, banner *domain.Banner) error {
	return r.db.WithContext(ctx).Save(banner).Error
}

func (r *bannerRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.Banner{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *bannerRepository) ListActive(ctx context.Context) ([]domain.Banner, error) {
	now := time.Now()
	var banners []domain.Banner
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL AND is_active = ?", true).
		Where("start_date IS NULL OR start_date <= ?", now).
		Where("end_date IS NULL OR end_date >= ?", now).
		Order("sort_order ASC").
		Find(&banners).Error
	return banners, err
}

func (r *bannerRepository) List(ctx context.Context) ([]domain.Banner, error) {
	var banners []domain.Banner
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("sort_order ASC").
		Find(&banners).Error
	return banners, err
}
