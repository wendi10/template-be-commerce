package domain

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID  `json:"id"               gorm:"primaryKey;type:char(36)"`
	UserID    uuid.UUID  `json:"user_id"          gorm:"type:char(36);not null;index"`
	FirstName string     `json:"first_name"       gorm:"type:varchar(100);not null;default:''"`
	LastName  string     `json:"last_name"        gorm:"type:varchar(100);not null;default:''"`
	Phone     string     `json:"phone"            gorm:"type:varchar(20)"`
	Avatar    string     `json:"avatar,omitempty" gorm:"type:text"`
	CreatedAt time.Time  `json:"created_at"       gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at"       gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Associations
	User      *User     `json:"user,omitempty"      gorm:"foreignKey:UserID"`
	Addresses []Address `json:"addresses,omitempty" gorm:"foreignKey:CustomerID"`
}

func (Customer) TableName() string { return "customers" }

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  string `json:"last_name"  validate:"omitempty,min=2,max=100"`
	Phone     string `json:"phone"      validate:"omitempty"`
	Avatar    string `json:"avatar"     validate:"omitempty,url"`
}

type Address struct {
	ID            uuid.UUID  `json:"id"                     gorm:"primaryKey;type:char(36)"`
	CustomerID    uuid.UUID  `json:"customer_id"            gorm:"type:char(36);not null;index"`
	Label         string     `json:"label"                  gorm:"type:varchar(100);default:'Home'"`
	RecipientName string     `json:"recipient_name"         gorm:"type:varchar(150);not null"`
	Phone         string     `json:"phone"                  gorm:"type:varchar(20);not null"`
	AddressLine1  string     `json:"address_line1"          gorm:"type:varchar(255);not null"`
	AddressLine2  string     `json:"address_line2,omitempty" gorm:"type:varchar(255)"`
	City          string     `json:"city"                   gorm:"type:varchar(100);not null"`
	Province      string     `json:"province"               gorm:"type:varchar(100);not null"`
	PostalCode    string     `json:"postal_code"            gorm:"type:varchar(10);not null"`
	IsDefault     bool       `json:"is_default"             gorm:"default:false;not null"`
	CreatedAt     time.Time  `json:"created_at"             gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at"             gorm:"autoUpdateTime"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"   gorm:"index"`
}

func (Address) TableName() string { return "addresses" }

type CreateAddressRequest struct {
	Label         string `json:"label"          validate:"required,max=50"`
	RecipientName string `json:"recipient_name" validate:"required,max=100"`
	Phone         string `json:"phone"          validate:"required"`
	AddressLine1  string `json:"address_line1"  validate:"required,max=255"`
	AddressLine2  string `json:"address_line2"  validate:"omitempty,max=255"`
	City          string `json:"city"           validate:"required,max=100"`
	Province      string `json:"province"       validate:"required,max=100"`
	PostalCode    string `json:"postal_code"    validate:"required,max=10"`
	IsDefault     bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	Label         string `json:"label"          validate:"omitempty,max=50"`
	RecipientName string `json:"recipient_name" validate:"omitempty,max=100"`
	Phone         string `json:"phone"          validate:"omitempty"`
	AddressLine1  string `json:"address_line1"  validate:"omitempty,max=255"`
	AddressLine2  string `json:"address_line2"  validate:"omitempty,max=255"`
	City          string `json:"city"           validate:"omitempty,max=100"`
	Province      string `json:"province"       validate:"omitempty,max=100"`
	PostalCode    string `json:"postal_code"    validate:"omitempty,max=10"`
	IsDefault     *bool  `json:"is_default"`
}
