package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/repository"
	"github.com/template-be-commerce/pkg/apperrors"
)

type PromoUseCase interface {
	Create(ctx context.Context, req domain.CreatePromoRequest) (*domain.PromoCode, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdatePromoRequest) (*domain.PromoCode, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PromoCode, error)
	List(ctx context.Context) ([]domain.PromoCode, error)
	ValidateAndCalculate(ctx context.Context, code string, subtotal decimal.Decimal) (*domain.ValidatePromoResponse, error)
}

type promoUseCase struct {
	promoRepo repository.PromoRepository
}

func NewPromoUseCase(promoRepo repository.PromoRepository) PromoUseCase {
	return &promoUseCase{promoRepo: promoRepo}
}

func (uc *promoUseCase) Create(ctx context.Context, req domain.CreatePromoRequest) (*domain.PromoCode, error) {
	promo := &domain.PromoCode{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Value:       req.Value,
		MinPurchase: req.MinPurchase,
		MaxDiscount: req.MaxDiscount,
		UsageLimit:  req.UsageLimit,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		IsActive:    req.IsActive,
	}
	if err := uc.promoRepo.Create(ctx, promo); err != nil {
		return nil, err
	}
	return promo, nil
}

func (uc *promoUseCase) Update(ctx context.Context, id uuid.UUID, req domain.UpdatePromoRequest) (*domain.PromoCode, error) {
	promo, err := uc.promoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Code != "" {
		promo.Code = req.Code
	}
	if req.Name != "" {
		promo.Name = req.Name
	}
	if req.Description != "" {
		promo.Description = req.Description
	}
	if req.Type != "" {
		promo.Type = req.Type
	}
	if req.Value != nil {
		promo.Value = *req.Value
	}
	if req.MinPurchase != nil {
		promo.MinPurchase = *req.MinPurchase
	}
	if req.MaxDiscount != nil {
		promo.MaxDiscount = *req.MaxDiscount
	}
	if req.UsageLimit != nil {
		promo.UsageLimit = *req.UsageLimit
	}
	if req.StartDate != nil {
		promo.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		promo.EndDate = *req.EndDate
	}
	if req.IsActive != nil {
		promo.IsActive = *req.IsActive
	}

	if err := uc.promoRepo.Update(ctx, promo); err != nil {
		return nil, err
	}
	return promo, nil
}

func (uc *promoUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.promoRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.promoRepo.SoftDelete(ctx, id)
}

func (uc *promoUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.PromoCode, error) {
	return uc.promoRepo.GetByID(ctx, id)
}

func (uc *promoUseCase) List(ctx context.Context) ([]domain.PromoCode, error) {
	return uc.promoRepo.List(ctx)
}

func (uc *promoUseCase) ValidateAndCalculate(ctx context.Context, code string, subtotal decimal.Decimal) (*domain.ValidatePromoResponse, error) {
	promo, err := uc.promoRepo.GetByCode(ctx, code)
	if err != nil {
		return &domain.ValidatePromoResponse{IsValid: false, Message: "invalid promo code"}, nil
	}

	now := time.Now()
	if !promo.IsActive {
		return &domain.ValidatePromoResponse{IsValid: false, Message: "promo code is not active"}, apperrors.ErrInvalidPromo
	}
	if now.Before(promo.StartDate) || now.After(promo.EndDate) {
		return &domain.ValidatePromoResponse{IsValid: false, Message: "promo code is expired"}, apperrors.ErrPromoExpired
	}
	if promo.UsageLimit > 0 && promo.UsedCount >= promo.UsageLimit {
		return &domain.ValidatePromoResponse{IsValid: false, Message: "promo code usage limit reached"}, apperrors.ErrPromoUsageLimitReached
	}
	if subtotal.LessThan(promo.MinPurchase) {
		return &domain.ValidatePromoResponse{
			IsValid: false,
			Message: "minimum purchase not met for this promo",
		}, apperrors.BadRequest("minimum purchase for promo is " + promo.MinPurchase.String())
	}

	var discount decimal.Decimal
	if promo.Type == domain.PromoTypePercentage {
		discount = subtotal.Mul(promo.Value).Div(decimal.NewFromInt(100))
		if promo.MaxDiscount.IsPositive() && discount.GreaterThan(promo.MaxDiscount) {
			discount = promo.MaxDiscount
		}
	} else {
		discount = promo.Value
		if discount.GreaterThan(subtotal) {
			discount = subtotal
		}
	}

	return &domain.ValidatePromoResponse{
		PromoCode:      promo,
		DiscountAmount: discount,
		IsValid:        true,
	}, nil
}
