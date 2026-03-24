package mysql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/pkg/apperrors"
)

// ─── Customer Repository ──────────────────────────────────────────────────────

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *customerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	if customer.ID == uuid.Nil {
		customer.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	var c domain.Customer
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *customerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Customer, error) {
	var c domain.Customer
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	return r.db.WithContext(ctx).
		Model(customer).
		Select("first_name", "last_name", "phone", "avatar").
		Updates(customer).Error
}

func (r *customerRepository) List(ctx context.Context, page, limit int) ([]domain.Customer, int64, error) {
	var customers []domain.Customer
	var total int64

	base := r.db.WithContext(ctx).Model(&domain.Customer{}).Where("deleted_at IS NULL")

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := base.
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&customers).Error

	return customers, total, err
}

// ─── Address Repository ───────────────────────────────────────────────────────

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *addressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, address *domain.Address) error {
	if address.ID == uuid.Nil {
		address.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *addressRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	var a domain.Address
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *addressRepository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Address, error) {
	var addresses []domain.Address
	err := r.db.WithContext(ctx).
		Where("customer_id = ? AND deleted_at IS NULL", customerID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

func (r *addressRepository) Update(ctx context.Context, address *domain.Address) error {
	return r.db.WithContext(ctx).Save(address).Error
}

func (r *addressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.Address{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *addressRepository) SetDefault(ctx context.Context, customerID, addressID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all defaults for this customer
		if err := tx.Model(&domain.Address{}).
			Where("customer_id = ?", customerID).
			Update("is_default", false).Error; err != nil {
			return err
		}
		// Set the target address as default
		return tx.Model(&domain.Address{}).
			Where("id = ?", addressID).
			Update("is_default", true).Error
	})
}
