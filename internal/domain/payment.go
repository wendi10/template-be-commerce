package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusSuccess  PaymentStatus = "success"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusExpired  PaymentStatus = "expired"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type PaymentProvider string

const (
	PaymentProviderDoku PaymentProvider = "doku"
)

type Payment struct {
	ID            uuid.UUID       `json:"id"                     gorm:"primaryKey;type:char(36)"`
	OrderID       uuid.UUID       `json:"order_id"               gorm:"type:char(36);not null;uniqueIndex"`
	Provider      PaymentProvider `json:"provider"               gorm:"type:varchar(50);not null"`
	// provider_ref: transaction/reference ID from the payment gateway
	TransactionID string          `json:"transaction_id,omitempty" gorm:"column:provider_ref;type:varchar(255);index"`
	// method: e.g. virtual_account, credit_card, qris
	PaymentMethod string          `json:"payment_method"         gorm:"column:method;type:varchar(50)"`
	Status        PaymentStatus   `json:"status"                 gorm:"type:enum('pending','success','failed','expired','refunded');default:'pending';not null"`
	Amount        decimal.Decimal `json:"amount"                 gorm:"type:decimal(15,2);not null"`
	Currency      string          `json:"currency"               gorm:"type:char(3);default:'IDR';not null"`
	// Payment URL returned by the provider to redirect the customer
	PaymentURL    string          `json:"payment_url,omitempty"  gorm:"column:payment_url;type:text"`
	// Raw JSON payload sent to the provider
	Payload       string          `json:"payload,omitempty"      gorm:"type:json"`
	// Raw JSON callback received from the provider
	CallbackData  string          `json:"callback_data,omitempty" gorm:"column:callback_payload;type:json"`
	PaidAt        *time.Time      `json:"paid_at,omitempty"`
	ExpiredAt     *time.Time      `json:"expired_at,omitempty"`
	CreatedAt     time.Time       `json:"created_at"             gorm:"autoCreateTime"`
	UpdatedAt     time.Time       `json:"updated_at"             gorm:"autoUpdateTime"`

	// Associations
	Order *Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

func (Payment) TableName() string { return "payments" }

type CreatePaymentRequest struct {
	OrderID       uuid.UUID       `json:"order_id"        validate:"required"`
	PaymentMethod string          `json:"payment_method"  validate:"required"`
	Provider      PaymentProvider `json:"provider"        validate:"required,oneof=doku"`
}

type PaymentCallbackData struct {
	TransactionID string          `json:"transaction_id"`
	OrderID       string          `json:"order_id"`
	Amount        decimal.Decimal `json:"amount"`
	Status        PaymentStatus   `json:"status"`
	Provider      PaymentProvider `json:"provider"`
	RawPayload    string          `json:"raw_payload"`
}

// CheckoutRequest combines order creation with payment initiation
type CheckoutRequest struct {
	AddressID     uuid.UUID       `json:"address_id"      validate:"required"`
	PromoCode     string          `json:"promo_code"`
	Notes         string          `json:"notes"           validate:"omitempty,max=500"`
	PaymentMethod string          `json:"payment_method"  validate:"required"`
	Provider      PaymentProvider `json:"provider"        validate:"required,oneof=doku"`
}

type CheckoutResponse struct {
	Order      Order   `json:"order"`
	Payment    Payment `json:"payment"`
	PaymentURL string  `json:"payment_url"`
}
