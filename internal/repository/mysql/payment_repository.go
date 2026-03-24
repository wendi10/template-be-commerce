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

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *paymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	if payment.ID == uuid.Nil {
		payment.ID = uuid.New()
	}
	return r.db.WithContext(ctx).
		Omit("Order").
		Create(payment).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

// GetByTransactionID maps to the provider_ref column (GORM column tag on TransactionID).
func (r *paymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.WithContext(ctx).
		Where("provider_ref = ?", transactionID).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	payment.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).
		Omit("Order").
		Save(payment).Error
}
