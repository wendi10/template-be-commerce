package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/repository"
	"github.com/template-be-commerce/pkg/apperrors"
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, req domain.CreateOrderRequest) (*domain.Order, error)
	GetOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, isAdmin bool) (*domain.Order, error)
	ListOrders(ctx context.Context, userID uuid.UUID, filter domain.OrderListFilter, isAdmin bool) (*domain.OrderListResponse, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, req domain.UpdateOrderStatusRequest) (*domain.Order, error)
	CancelOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) error
	GetSalesReport(ctx context.Context, startDate, endDate time.Time) (*domain.SalesReport, error)
}

type orderUseCase struct {
	orderRepo    repository.OrderRepository
	cartRepo     repository.CartRepository
	customerRepo repository.CustomerRepository
	productRepo  repository.ProductRepository
	promoRepo    repository.PromoRepository
	addressRepo  repository.AddressRepository
	promoUC      PromoUseCase
}

func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	customerRepo repository.CustomerRepository,
	productRepo repository.ProductRepository,
	promoRepo repository.PromoRepository,
	addressRepo repository.AddressRepository,
	promoUC PromoUseCase,
) OrderUseCase {
	return &orderUseCase{
		orderRepo:    orderRepo,
		cartRepo:     cartRepo,
		customerRepo: customerRepo,
		productRepo:  productRepo,
		promoRepo:    promoRepo,
		addressRepo:  addressRepo,
		promoUC:      promoUC,
	}
}

func (uc *orderUseCase) CreateOrder(ctx context.Context, userID uuid.UUID, req domain.CreateOrderRequest) (*domain.Order, error) {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Validate address ownership
	address, err := uc.addressRepo.GetByID(ctx, req.AddressID)
	if err != nil {
		return nil, err
	}
	if address.CustomerID != customer.ID {
		return nil, apperrors.ErrForbidden
	}

	// Load cart
	cartItems, err := uc.cartRepo.ListByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, apperrors.ErrCartEmpty
	}

	// Build order items and check stock
	var orderItems []domain.OrderItem
	subTotal := decimal.Zero

	for _, cartItem := range cartItems {
		product, err := uc.productRepo.GetByID(ctx, cartItem.ProductID)
		if err != nil {
			return nil, err
		}
		if !product.IsActive {
			return nil, apperrors.BadRequest(fmt.Sprintf("product '%s' is no longer available", product.Name))
		}
		if product.Stock < cartItem.Quantity {
			return nil, apperrors.BadRequest(fmt.Sprintf("insufficient stock for '%s'", product.Name))
		}

		lineTotal := product.Price.Mul(decimal.NewFromInt(int64(cartItem.Quantity)))
		subTotal = subTotal.Add(lineTotal)

		orderItems = append(orderItems, domain.OrderItem{
			ProductID:    product.ID,
			ProductName:  product.Name,
			ProductPrice: product.Price,
			Quantity:     cartItem.Quantity,
			TotalPrice:   lineTotal,
		})
	}

	// Apply promo
	var discountAmount decimal.Decimal
	var promoCodeID *uuid.UUID

	if req.PromoCode != "" {
		promoResult, err := uc.promoUC.ValidateAndCalculate(ctx, req.PromoCode, subTotal)
		if err != nil {
			return nil, err
		}
		if promoResult.IsValid {
			discountAmount = promoResult.DiscountAmount
			promoCodeID = &promoResult.PromoCode.ID
		}
	}

	shippingCost := decimal.NewFromFloat(15000) // flat rate, can be replaced with shipping API
	totalAmount := subTotal.Sub(discountAmount).Add(shippingCost)

	order := &domain.Order{
		CustomerID:     customer.ID,
		AddressID:      req.AddressID,
		PromoCodeID:    promoCodeID,
		OrderNumber:    generateOrderNumber(),
		Status:         domain.OrderStatusPending,
		SubTotal:       subTotal,
		DiscountAmount: discountAmount,
		ShippingCost:   shippingCost,
		TotalAmount:    totalAmount,
		Notes:          req.Notes,
	}

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Assign order ID to items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}

	if err := uc.orderRepo.CreateItems(ctx, orderItems); err != nil {
		return nil, err
	}
	order.Items = orderItems

	// Decrement stock
	for _, item := range cartItems {
		if err := uc.productRepo.DecrementStock(ctx, item.ProductID, item.Quantity); err != nil {
			return nil, err
		}
	}

	// Increment promo usage
	if promoCodeID != nil {
		_ = uc.promoRepo.IncrementUsage(ctx, *promoCodeID)
	}

	// Clear cart
	_ = uc.cartRepo.ClearCart(ctx, customer.ID)

	return order, nil
}

func (uc *orderUseCase) GetOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, isAdmin bool) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		customer, err := uc.customerRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if order.CustomerID != customer.ID {
			return nil, apperrors.ErrForbidden
		}
	}

	return order, nil
}

func (uc *orderUseCase) ListOrders(ctx context.Context, userID uuid.UUID, filter domain.OrderListFilter, isAdmin bool) (*domain.OrderListResponse, error) {
	if !isAdmin {
		customer, err := uc.customerRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		filter.CustomerID = &customer.ID
	}

	orders, total, err := uc.orderRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if filter.Limit > 0 {
		totalPages = int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	}

	return &domain.OrderListResponse{
		Orders:     orders,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (uc *orderUseCase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, req domain.UpdateOrderStatusRequest) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Validate transition
	allowed := domain.ValidTransitions[order.Status]
	valid := false
	for _, s := range allowed {
		if s == req.Status {
			valid = true
			break
		}
	}
	if !valid {
		return nil, apperrors.ErrInvalidOrderTransition
	}

	if err := uc.orderRepo.UpdateStatus(ctx, orderID, req.Status); err != nil {
		return nil, err
	}
	order.Status = req.Status
	return order, nil
}

func (uc *orderUseCase) CancelOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) error {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.CustomerID != customer.ID {
		return apperrors.ErrForbidden
	}

	allowed := domain.ValidTransitions[order.Status]
	canCancel := false
	for _, s := range allowed {
		if s == domain.OrderStatusCancelled {
			canCancel = true
			break
		}
	}
	if !canCancel {
		return apperrors.ErrInvalidOrderTransition
	}

	// Restore stock
	for _, item := range order.Items {
		_ = uc.productRepo.IncrementStock(ctx, item.ProductID, item.Quantity)
	}

	return uc.orderRepo.UpdateStatus(ctx, orderID, domain.OrderStatusCancelled)
}

func (uc *orderUseCase) GetSalesReport(ctx context.Context, startDate, endDate time.Time) (*domain.SalesReport, error) {
	return uc.orderRepo.SalesReport(ctx, startDate, endDate)
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixMilli())
}
