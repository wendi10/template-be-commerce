package domain

import (
	"time"

	"github.com/google/uuid"
)

type Banner struct {
	ID        uuid.UUID  `json:"id"                   gorm:"primaryKey;type:char(36)"`
	Title     string     `json:"title"                gorm:"type:varchar(255);not null"`
	Subtitle  string     `json:"subtitle,omitempty"   gorm:"type:varchar(255)"`
	ImageURL  string     `json:"image_url"            gorm:"column:image_url;type:text;not null"`
	LinkURL   string     `json:"link_url,omitempty"   gorm:"column:link_url;type:text"`
	IsActive  bool       `json:"is_active"            gorm:"default:true;not null"`
	SortOrder int        `json:"sort_order"           gorm:"default:0;not null"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"           gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at"           gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

func (Banner) TableName() string { return "banners" }

type CreateBannerRequest struct {
	Title     string     `json:"title"     validate:"required,max=255"`
	Subtitle  string     `json:"subtitle"  validate:"omitempty,max=500"`
	ImageURL  string     `json:"image_url" validate:"required,url"`
	LinkURL   string     `json:"link_url"  validate:"omitempty,url"`
	IsActive  bool       `json:"is_active"`
	SortOrder int        `json:"sort_order"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type UpdateBannerRequest struct {
	Title     string     `json:"title"     validate:"omitempty,max=255"`
	Subtitle  string     `json:"subtitle"  validate:"omitempty,max=500"`
	ImageURL  string     `json:"image_url" validate:"omitempty,url"`
	LinkURL   string     `json:"link_url"  validate:"omitempty,url"`
	IsActive  *bool      `json:"is_active"`
	SortOrder *int       `json:"sort_order"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}
