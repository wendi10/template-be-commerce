package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/payment"
	"github.com/template-be-commerce/internal/repository"
	"github.com/template-be-commerce/pkg/apperrors"
	"github.com/template-be-commerce/pkg/logger"
	"go.uber.org/zap"
)

type PaymentUseCase interface {
	CreatePayment(ctx context.Context, userID uuid.UUID, req domain.CreatePaymentRequest) (*domain.CheckoutResponse, error)
	HandleCallback(ctx context.Context, provider domain.PaymentProvider, payload []byte, headers map[string]string) error
	GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error)
}

type paymentUseCase struct {
	paymentRepo  repository.PaymentRepository
	orderRepo    repository.OrderRepository
	customerRepo repository.CustomerRepository
	addressRepo  repository.AddressRepository
	gateways     map[domain.PaymentProvider]payment.Gateway
	callbackURL  string
}

func NewPaymentUseCase(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	customerRepo repository.CustomerRepository,
	addressRepo repository.AddressRepository,
	gateways map[domain.PaymentProvider]payment.Gateway,
	callbackURL string,
) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo:  paymentRepo,
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		addressRepo:  addressRepo,
		gateways:     gateways,
		callbackURL:  callbackURL,
	}
}

func (uc *paymentUseCase) CreatePayment(ctx context.Context, userID uuid.UUID, req domain.CreatePaymentRequest) (*domain.CheckoutResponse, error) {
	customer, err := uc.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	order, err := uc.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	if order.CustomerID != customer.ID {
		return nil, apperrors.ErrForbidden
	}

	if order.Status != domain.OrderStatusPending {
		return nil, apperrors.BadRequest("order is not in a payable state")
	}

	gateway, ok := uc.gateways[req.Provider]
	if !ok {
		return nil, apperrors.BadRequest("payment provider not supported")
	}

	expiredAt := time.Now().Add(60 * time.Minute).Unix()
	txReq := payment.CreateTransactionRequest{
		OrderID:       order.ID.String(),
		OrderNumber:   order.OrderNumber,
		Amount:        order.TotalAmount.String(),
		PaymentMethod: req.PaymentMethod,
		CustomerName:  customer.FirstName + " " + customer.LastName,
		CustomerPhone: customer.Phone,
		Description:   "Payment for order " + order.OrderNumber,
		CallbackURL:   uc.callbackURL,
		ExpiredAt:     expiredAt,
	}

	// Get customer email from user
	txResp, err := gateway.CreateTransaction(ctx, txReq)
	if err != nil {
		logger.Error("payment gateway error", zap.Error(err))
		return nil, apperrors.InternalServer("failed to create payment transaction", err)
	}

	expiredTime := time.Unix(txResp.ExpiredAt, 0)
	pay := &domain.Payment{
		OrderID:       order.ID,
		PaymentMethod: req.PaymentMethod,
		Provider:      req.Provider,
		Amount:        order.TotalAmount,
		Status:        domain.PaymentStatusPending,
		TransactionID: txResp.TransactionID,
		PaymentURL:    txResp.PaymentURL,
		CallbackData:  txResp.RawResponse,
		ExpiredAt:     &expiredTime,
	}

	if err := uc.paymentRepo.Create(ctx, pay); err != nil {
		return nil, err
	}

	// Update order to waiting_payment
	_ = uc.orderRepo.UpdateStatus(ctx, order.ID, domain.OrderStatusWaitingPayment)
	order.Status = domain.OrderStatusWaitingPayment

	return &domain.CheckoutResponse{
		Order:      *order,
		Payment:    *pay,
		PaymentURL: txResp.PaymentURL,
	}, nil
}

func (uc *paymentUseCase) HandleCallback(ctx context.Context, provider domain.PaymentProvider, payload []byte, headers map[string]string) error {
	gateway, ok := uc.gateways[provider]
	if !ok {
		return apperrors.BadRequest("unknown payment provider")
	}

	result, err := gateway.HandleCallback(ctx, payload, headers)
	if err != nil {
		logger.Error("payment callback parse error", zap.Error(err), zap.String("provider", string(provider)))
		return err
	}

	pay, err := uc.paymentRepo.GetByTransactionID(ctx, result.TransactionID)
	if err != nil {
		// Try to find by order number (DOKU sends invoice_number)
		order, oErr := uc.orderRepo.GetByOrderNumber(ctx, result.OrderID)
		if oErr != nil {
			logger.Error("payment callback: order not found",
				zap.String("transaction_id", result.TransactionID),
				zap.String("order_id", result.OrderID),
			)
			return apperrors.NotFound("payment transaction not found")
		}
		pay, err = uc.paymentRepo.GetByOrderID(ctx, order.ID)
		if err != nil {
			return apperrors.NotFound("payment record not found")
		}
	}

	// Idempotency: skip if already processed
	if pay.Status == domain.PaymentStatusSuccess {
		return nil
	}

	pay.Status = result.Status
	pay.CallbackData = result.RawPayload
	if result.Status == domain.PaymentStatusSuccess {
		now := time.Now()
		pay.PaidAt = &now
	}

	if err := uc.paymentRepo.Update(ctx, pay); err != nil {
		return err
	}

	// Update order status based on payment result
	var newOrderStatus domain.OrderStatus
	switch result.Status {
	case domain.PaymentStatusSuccess:
		newOrderStatus = domain.OrderStatusPaid
	case domain.PaymentStatusFailed, domain.PaymentStatusExpired:
		newOrderStatus = domain.OrderStatusCancelled
	default:
		return nil
	}

	if err := uc.orderRepo.UpdateStatus(ctx, pay.OrderID, newOrderStatus); err != nil {
		logger.Error("failed to update order status after payment callback",
			zap.String("order_id", pay.OrderID.String()),
			zap.Error(err),
		)
	}

	return nil
}

func (uc *paymentUseCase) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	return uc.paymentRepo.GetByOrderID(ctx, orderID)
}
