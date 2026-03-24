package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CartItem struct {
	ID         uuid.UUID `json:"id"          gorm:"primaryKey;type:char(36)"`
	CustomerID uuid.UUID `json:"customer_id" gorm:"type:char(36);not null;index"`
	ProductID  uuid.UUID `json:"product_id"  gorm:"type:char(36);not null"`
	Quantity   int       `json:"quantity"    gorm:"default:1;not null"`
	CreatedAt  time.Time `json:"created_at"  gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at"  gorm:"autoUpdateTime"`

	// Associations & computed fields (not stored)
	Product  *Product        `json:"product,omitempty"   gorm:"foreignKey:ProductID"`
	SubTotal decimal.Decimal `json:"sub_total,omitempty" gorm:"-"`
}

func (CartItem) TableName() string { return "cart_items" }

type CartSummary struct {
	Items          []CartItem      `json:"items"`
	TotalItems     int             `json:"total_items"`
	SubTotal       decimal.Decimal `json:"sub_total"`
	PromoCode      *PromoCode      `json:"promo_code,omitempty"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	Total          decimal.Decimal `json:"total"`
}

type AddToCartRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity"   validate:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

type ApplyPromoToCartRequest struct {
	PromoCode string `json:"promo_code" validate:"required"`
}
