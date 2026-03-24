package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderStatus string

const (
	OrderStatusPending        OrderStatus = "pending"
	OrderStatusWaitingPayment OrderStatus = "waiting_payment"
	OrderStatusPaid           OrderStatus = "paid"
	OrderStatusProcessing     OrderStatus = "processing"
	OrderStatusShipped        OrderStatus = "shipped"
	OrderStatusCompleted      OrderStatus = "completed"
	OrderStatusCancelled      OrderStatus = "cancelled"
)

// ValidTransitions defines allowed status transitions
var ValidTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusPending:        {OrderStatusWaitingPayment, OrderStatusCancelled},
	OrderStatusWaitingPayment: {OrderStatusPaid, OrderStatusCancelled},
	OrderStatusPaid:           {OrderStatusProcessing, OrderStatusCancelled},
	OrderStatusProcessing:     {OrderStatusShipped},
	OrderStatusShipped:        {OrderStatusCompleted},
	OrderStatusCompleted:      {},
	OrderStatusCancelled:      {},
}

type Order struct {
	ID             uuid.UUID       `json:"id"               gorm:"primaryKey;type:char(36)"`
	CustomerID     uuid.UUID       `json:"customer_id"      gorm:"type:char(36);not null;index"`
	AddressID      uuid.UUID       `json:"address_id"       gorm:"type:char(36);not null"`
	PromoCodeID    *uuid.UUID      `json:"promo_code_id,omitempty" gorm:"type:char(36);index"`
	OrderNumber    string          `json:"order_number"     gorm:"type:varchar(50);not null;uniqueIndex"`
	Status         OrderStatus     `json:"status"           gorm:"type:enum('pending','waiting_payment','paid','processing','shipped','completed','cancelled');default:'pending';not null"`
	// Column "subtotal" in DB — GORM would map SubTotal → sub_total, so explicit column name needed
	SubTotal       decimal.Decimal `json:"sub_total"        gorm:"column:subtotal;type:decimal(15,2);not null;default:0.00"`
	DiscountAmount decimal.Decimal `json:"discount_amount"  gorm:"type:decimal(15,2);not null;default:0.00"`
	ShippingCost   decimal.Decimal `json:"shipping_cost"    gorm:"type:decimal(15,2);not null;default:0.00"`
	TotalAmount    decimal.Decimal `json:"total_amount"     gorm:"type:decimal(15,2);not null;default:0.00"`
	Notes          string          `json:"notes,omitempty"  gorm:"type:text"`
	CreatedAt      time.Time       `json:"created_at"       gorm:"autoCreateTime"`
	UpdatedAt      time.Time       `json:"updated_at"       gorm:"autoUpdateTime"`

	// Associations
	Customer  *Customer   `json:"customer,omitempty"   gorm:"foreignKey:CustomerID"`
	Address   *Address    `json:"address,omitempty"    gorm:"foreignKey:AddressID"`
	PromoCode *PromoCode  `json:"promo_code,omitempty" gorm:"foreignKey:PromoCodeID"`
	Items     []OrderItem `json:"items,omitempty"      gorm:"foreignKey:OrderID"`
	Payment   *Payment    `json:"payment,omitempty"    gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string { return "orders" }

type OrderItem struct {
	ID           uuid.UUID       `json:"id"            gorm:"primaryKey;type:char(36)"`
	OrderID      uuid.UUID       `json:"order_id"      gorm:"type:char(36);not null;index"`
	ProductID    uuid.UUID       `json:"product_id"    gorm:"type:char(36);not null;index"`
	ProductName  string          `json:"product_name"  gorm:"type:varchar(255);not null"`
	ProductSlug  string          `json:"product_slug"  gorm:"type:varchar(300);not null"`
	// "unit_price" in DB — GORM would map ProductPrice → product_price
	ProductPrice decimal.Decimal `json:"product_price" gorm:"column:unit_price;type:decimal(15,2);not null"`
	Quantity     int             `json:"quantity"      gorm:"default:1;not null"`
	// "subtotal" in DB — GORM would map TotalPrice → total_price
	TotalPrice   decimal.Decimal `json:"total_price"   gorm:"column:subtotal;type:decimal(15,2);not null"`
	CreatedAt    time.Time       `json:"created_at"    gorm:"autoCreateTime"`
}

func (OrderItem) TableName() string { return "order_items" }

type CreateOrderRequest struct {
	AddressID uuid.UUID `json:"address_id" validate:"required"`
	PromoCode string    `json:"promo_code"`
	Notes     string    `json:"notes"      validate:"omitempty,max=500"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required,oneof=pending waiting_payment paid processing shipped completed cancelled"`
	Notes  string      `json:"notes"`
}

type OrderListFilter struct {
	CustomerID *uuid.UUID
	Status     *OrderStatus
	StartDate  *time.Time
	EndDate    *time.Time
	Page       int
	Limit      int
}

type OrderListResponse struct {
	Orders     []Order `json:"orders"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalPages int     `json:"total_pages"`
}

type SalesReport struct {
	TotalOrders     int64           `json:"total_orders"`
	TotalRevenue    decimal.Decimal `json:"total_revenue"`
	AverageOrder    decimal.Decimal `json:"average_order"`
	PendingOrders   int64           `json:"pending_orders"`
	CompletedOrders int64           `json:"completed_orders"`
	CancelledOrders int64           `json:"cancelled_orders"`
	TopProducts     []TopProduct    `json:"top_products"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         time.Time       `json:"end_date"`
}

type TopProduct struct {
	ProductID   uuid.UUID       `json:"product_id"`
	ProductName string          `json:"product_name"`
	TotalSold   int64           `json:"total_sold"`
	Revenue     decimal.Decimal `json:"revenue"`
}
