package mysql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/pkg/apperrors"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *cartRepository {
	return &cartRepository{db: db}
}

// AddItem inserts a cart item or increments quantity if it already exists
// (uses MySQL ON DUPLICATE KEY UPDATE via GORM's upsert).
func (r *cartRepository) AddItem(ctx context.Context, item *domain.CartItem) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "customer_id"}, {Name: "product_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"quantity":   gorm.Expr("cart_items.quantity + VALUES(quantity)"),
				"updated_at": gorm.Expr("NOW()"),
			}),
		}).
		Create(item).Error
}

func (r *cartRepository) GetItem(ctx context.Context, customerID, productID uuid.UUID) (*domain.CartItem, error) {
	var item domain.CartItem
	err := r.db.WithContext(ctx).
		Where("customer_id = ? AND product_id = ?", customerID, productID).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) GetItemByID(ctx context.Context, id uuid.UUID) (*domain.CartItem, error) {
	var item domain.CartItem
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) UpdateItem(ctx context.Context, item *domain.CartItem) error {
	return r.db.WithContext(ctx).
		Model(item).
		Update("quantity", item.Quantity).Error
}

func (r *cartRepository) RemoveItem(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.CartItem{}).Error
}

func (r *cartRepository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.CartItem, error) {
	var items []domain.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL")
		}).
		Preload("Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true).Limit(1)
		}).
		Where("customer_id = ?", customerID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *cartRepository) ClearCart(ctx context.Context, customerID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("customer_id = ?", customerID).
		Delete(&domain.CartItem{}).Error
}
