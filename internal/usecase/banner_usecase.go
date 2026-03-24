package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/repository"
)

type BannerUseCase interface {
	Create(ctx context.Context, req domain.CreateBannerRequest) (*domain.Banner, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdateBannerRequest) (*domain.Banner, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Banner, error)
	ListActive(ctx context.Context) ([]domain.Banner, error)
	ListAll(ctx context.Context) ([]domain.Banner, error)
}

type bannerUseCase struct {
	bannerRepo repository.BannerRepository
}

func NewBannerUseCase(bannerRepo repository.BannerRepository) BannerUseCase {
	return &bannerUseCase{bannerRepo: bannerRepo}
}

func (uc *bannerUseCase) Create(ctx context.Context, req domain.CreateBannerRequest) (*domain.Banner, error) {
	banner := &domain.Banner{
		Title:     req.Title,
		Subtitle:  req.Subtitle,
		ImageURL:  req.ImageURL,
		LinkURL:   req.LinkURL,
		IsActive:  req.IsActive,
		SortOrder: req.SortOrder,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	if err := uc.bannerRepo.Create(ctx, banner); err != nil {
		return nil, err
	}
	return banner, nil
}

func (uc *bannerUseCase) Update(ctx context.Context, id uuid.UUID, req domain.UpdateBannerRequest) (*domain.Banner, error) {
	banner, err := uc.bannerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		banner.Title = req.Title
	}
	if req.Subtitle != "" {
		banner.Subtitle = req.Subtitle
	}
	if req.ImageURL != "" {
		banner.ImageURL = req.ImageURL
	}
	if req.LinkURL != "" {
		banner.LinkURL = req.LinkURL
	}
	if req.IsActive != nil {
		banner.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		banner.SortOrder = *req.SortOrder
	}
	if req.StartDate != nil {
		banner.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		banner.EndDate = req.EndDate
	}

	if err := uc.bannerRepo.Update(ctx, banner); err != nil {
		return nil, err
	}
	return banner, nil
}

func (uc *bannerUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.bannerRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.bannerRepo.SoftDelete(ctx, id)
}

func (uc *bannerUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Banner, error) {
	return uc.bannerRepo.GetByID(ctx, id)
}

func (uc *bannerUseCase) ListActive(ctx context.Context) ([]domain.Banner, error) {
	return uc.bannerRepo.ListActive(ctx)
}

func (uc *bannerUseCase) ListAll(ctx context.Context) ([]domain.Banner, error) {
	return uc.bannerRepo.List(ctx)
}
