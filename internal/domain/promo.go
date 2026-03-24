package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PromoType string

const (
	PromoTypePercentage PromoType = "percentage"
	PromoTypeFixed      PromoType = "fixed"
)

type PromoCode struct {
	ID          uuid.UUID       `json:"id"                    gorm:"primaryKey;type:char(36)"`
	Code        string          `json:"code"                  gorm:"type:varchar(50);not null;uniqueIndex"`
	Name        string          `json:"name"                  gorm:"type:varchar(255);not null"`
	Description string          `json:"description,omitempty" gorm:"type:text"`
	// Column name differs from field name — mapped via gorm tag
	Type        PromoType       `json:"type"         gorm:"column:discount_type;type:enum('percentage','fixed');not null"`
	Value       decimal.Decimal `json:"value"        gorm:"column:discount_value;type:decimal(15,2);not null"`
	MinPurchase decimal.Decimal `json:"min_purchase" gorm:"type:decimal(15,2);not null;default:0.00"`
	MaxDiscount decimal.Decimal `json:"max_discount" gorm:"type:decimal(15,2);not null;default:0.00"`
	UsageLimit  int             `json:"usage_limit"  gorm:"default:0;not null"` // 0 = unlimited
	UsedCount   int             `json:"used_count"   gorm:"default:0;not null"`
	StartDate   time.Time       `json:"start_date"`
	EndDate     time.Time       `json:"end_date"`
	IsActive    bool            `json:"is_active"    gorm:"default:true;not null"`
	CreatedAt   time.Time       `json:"created_at"   gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at"   gorm:"autoUpdateTime"`
	DeletedAt   *time.Time      `json:"deleted_at,omitempty" gorm:"index"`
}

func (PromoCode) TableName() string { return "promo_codes" }

type CreatePromoRequest struct {
	Code        string          `json:"code"         validate:"required,max=50"`
	Name        string          `json:"name"         validate:"required,max=255"`
	Description string          `json:"description"`
	Type        PromoType       `json:"type"         validate:"required,oneof=percentage fixed"`
	Value       decimal.Decimal `json:"value"        validate:"required"`
	MinPurchase decimal.Decimal `json:"min_purchase"`
	MaxDiscount decimal.Decimal `json:"max_discount"`
	UsageLimit  int             `json:"usage_limit"  validate:"min=0"`
	StartDate   time.Time       `json:"start_date"   validate:"required"`
	EndDate     time.Time       `json:"end_date"     validate:"required"`
	IsActive    bool            `json:"is_active"`
}

type UpdatePromoRequest struct {
	Code        string           `json:"code"         validate:"omitempty,max=50"`
	Name        string           `json:"name"         validate:"omitempty,max=255"`
	Description string           `json:"description"`
	Type        PromoType        `json:"type"         validate:"omitempty,oneof=percentage fixed"`
	Value       *decimal.Decimal `json:"value"`
	MinPurchase *decimal.Decimal `json:"min_purchase"`
	MaxDiscount *decimal.Decimal `json:"max_discount"`
	UsageLimit  *int             `json:"usage_limit"  validate:"omitempty,min=0"`
	StartDate   *time.Time       `json:"start_date"`
	EndDate     *time.Time       `json:"end_date"`
	IsActive    *bool            `json:"is_active"`
}

type ValidatePromoResponse struct {
	PromoCode      *PromoCode      `json:"promo_code"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	IsValid        bool            `json:"is_valid"`
	Message        string          `json:"message,omitempty"`
}
