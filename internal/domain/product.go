package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Category struct {
	ID          uuid.UUID  `json:"id"                    gorm:"primaryKey;type:char(36)"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"   gorm:"type:char(36);index"`
	Name        string     `json:"name"                  gorm:"type:varchar(150);not null"`
	Slug        string     `json:"slug"                  gorm:"type:varchar(200);not null;uniqueIndex"`
	Description string     `json:"description,omitempty" gorm:"type:text"`
	IsActive    bool       `json:"is_active"             gorm:"default:true;not null"`
	CreatedAt   time.Time  `json:"created_at"            gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at"            gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"  gorm:"index"`

	// Associations
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

func (Category) TableName() string { return "categories" }

type Product struct {
	ID          uuid.UUID       `json:"id"                    gorm:"primaryKey;type:char(36)"`
	CategoryID  uuid.UUID       `json:"category_id"           gorm:"type:char(36);not null;index"`
	Name        string          `json:"name"                  gorm:"type:varchar(255);not null"`
	Slug        string          `json:"slug"                  gorm:"type:varchar(300);not null;uniqueIndex"`
	Description string          `json:"description"           gorm:"type:text"`
	Price       decimal.Decimal `json:"price"                 gorm:"type:decimal(15,2);not null;default:0.00"`
	Weight      decimal.Decimal `json:"weight"                gorm:"type:decimal(10,3);not null;default:0.000"`
	Stock       int             `json:"stock"                 gorm:"default:0;not null"`
	IsActive    bool            `json:"is_active"             gorm:"default:true;not null"`
	CreatedAt   time.Time       `json:"created_at"            gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at"            gorm:"autoUpdateTime"`
	DeletedAt   *time.Time      `json:"deleted_at,omitempty"  gorm:"index"`

	// Associations
	Category *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Images   []ProductImage `json:"images,omitempty"   gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string { return "products" }

type ProductImage struct {
	ID        uuid.UUID `json:"id"                  gorm:"primaryKey;type:char(36)"`
	ProductID uuid.UUID `json:"product_id"          gorm:"type:char(36);not null;index"`
	URL       string    `json:"url"                 gorm:"type:text;not null"`
	AltText   string    `json:"alt_text,omitempty"  gorm:"type:varchar(255)"`
	IsPrimary bool      `json:"is_primary"          gorm:"default:false;not null"`
	SortOrder int       `json:"sort_order"          gorm:"default:0;not null"`
	CreatedAt time.Time `json:"created_at"          gorm:"autoCreateTime"`
}

func (ProductImage) TableName() string { return "product_images" }

// --- Request/Response DTOs ---

type CreateProductRequest struct {
	CategoryID  uuid.UUID       `json:"category_id"  validate:"required"`
	Name        string          `json:"name"         validate:"required,min=2,max=255"`
	Slug        string          `json:"slug"         validate:"omitempty,max=255"`
	Description string          `json:"description"  validate:"required"`
	Price       decimal.Decimal `json:"price"        validate:"required"`
	Weight      decimal.Decimal `json:"weight"       validate:"required"`
	Stock       int             `json:"stock"        validate:"required,min=0"`
	IsActive    bool            `json:"is_active"`
}

type UpdateProductRequest struct {
	CategoryID  *uuid.UUID       `json:"category_id"`
	Name        string           `json:"name"        validate:"omitempty,min=2,max=255"`
	Slug        string           `json:"slug"        validate:"omitempty,max=255"`
	Description string           `json:"description"`
	Price       *decimal.Decimal `json:"price"`
	Weight      *decimal.Decimal `json:"weight"`
	Stock       *int             `json:"stock"       validate:"omitempty,min=0"`
	IsActive    *bool            `json:"is_active"`
}

type CreateCategoryRequest struct {
	Name        string     `json:"name"        validate:"required,min=2,max=100"`
	Slug        string     `json:"slug"        validate:"omitempty,max=100"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsActive    bool       `json:"is_active"`
}

type UpdateCategoryRequest struct {
	Name        string     `json:"name"        validate:"omitempty,min=2,max=100"`
	Slug        string     `json:"slug"        validate:"omitempty,max=100"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsActive    *bool      `json:"is_active"`
}

type AddProductImageRequest struct {
	URL       string `json:"url"        validate:"required,url"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}

type ProductListFilter struct {
	CategoryID *uuid.UUID
	Search     string
	MinPrice   *decimal.Decimal
	MaxPrice   *decimal.Decimal
	IsActive   *bool
	Page       int
	Limit      int
	SortBy     string
	SortOrder  string
}

type ProductListResponse struct {
	Products   []Product `json:"products"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}
