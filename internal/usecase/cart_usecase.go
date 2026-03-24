package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/repository"
	"github.com/template-be-commerce/pkg/apperrors"
)

type CartUseCase interface {
	AddToCart(ctx context.Context, userID uuid.UUID, req domain.AddToCartRequest) (*domain.CartItem, error)
	UpdateCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req domain.UpdateCartItemRequest) (*domain.CartItem, error)
	RemoveCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error
	GetCartSummary(ctx context.Context, userID uuid.UUID, promoCode string) (*domain.CartSummary, error)
}

type cartUseCase struct {
	cartRepo     repository.CartRepository
	customerRepo repository.CustomerRepository
	productRepo  repository.ProductRepository
	promoUC      PromoUseCase
}

func NewCartUseCase(
	cartRepo repository.CartRepository,
	customerRepo repository.CustomerRepository,
	productRepo repository.ProductRepository,
	promoUC PromoUseCase,
) CartUseCase {
	return &cartUseCase{
		cartRepo:     cartRepo,
		customerRepo: customerRepo,
		productRepo:  productRepo,
		promoUC:      promoUC,
	}
}

func (uc *cartUseCase) AddToCart(ctx context.Context, userID uuid.UUID, req domain.AddToCartRequest) (*domain.CartItem, error) {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	product, err := uc.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	if !product.IsActive {
		return nil, apperrors.BadRequest("product is not available")
	}

	if product.Stock < req.Quantity {
		return nil, apperrors.ErrInsufficientStock
	}

	// Check if item already in cart
	existing, err := uc.cartRepo.GetItem(ctx, customer.ID, req.ProductID)
	if err == nil && existing != nil {
		// Update quantity
		existing.Quantity += req.Quantity
		if product.Stock < existing.Quantity {
			return nil, apperrors.ErrInsufficientStock
		}
		if err := uc.cartRepo.UpdateItem(ctx, existing); err != nil {
			return nil, err
		}
		existing.Product = product
		return existing, nil
	}

	item := &domain.CartItem{
		CustomerID: customer.ID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
	}

	if err := uc.cartRepo.AddItem(ctx, item); err != nil {
		return nil, err
	}
	item.Product = product
	return item, nil
}

func (uc *cartUseCase) UpdateCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req domain.UpdateCartItemRequest) (*domain.CartItem, error) {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	item, err := uc.cartRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if item.CustomerID != customer.ID {
		return nil, apperrors.ErrForbidden
	}

	product, err := uc.productRepo.GetByID(ctx, item.ProductID)
	if err != nil {
		return nil, err
	}

	if product.Stock < req.Quantity {
		return nil, apperrors.ErrInsufficientStock
	}

	item.Quantity = req.Quantity
	if err := uc.cartRepo.UpdateItem(ctx, item); err != nil {
		return nil, err
	}
	item.Product = product
	return item, nil
}

func (uc *cartUseCase) RemoveCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	item, err := uc.cartRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return err
	}

	if item.CustomerID != customer.ID {
		return apperrors.ErrForbidden
	}

	return uc.cartRepo.RemoveItem(ctx, itemID)
}

func (uc *cartUseCase) GetCartSummary(ctx context.Context, userID uuid.UUID, promoCode string) (*domain.CartSummary, error) {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	items, err := uc.cartRepo.ListByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, err
	}

	subTotal := decimal.Zero
	for i, item := range items {
		if item.Product != nil {
			lineTotal := item.Product.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))
			items[i].SubTotal = lineTotal
			subTotal = subTotal.Add(lineTotal)
		}
	}

	summary := &domain.CartSummary{
		Items:      items,
		TotalItems: len(items),
		SubTotal:   subTotal,
		Total:      subTotal,
	}

	if promoCode != "" {
		promoResult, _ := uc.promoUC.ValidateAndCalculate(ctx, promoCode, subTotal)
		if promoResult != nil && promoResult.IsValid {
			summary.PromoCode = promoResult.PromoCode
			summary.DiscountAmount = promoResult.DiscountAmount
			summary.Total = subTotal.Sub(promoResult.DiscountAmount)
		}
	}

	return summary, nil
}
