package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/template-be-commerce/internal/domain"
)

// UserRepository handles persistence for users.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

// CustomerRepository handles persistence for customer profiles.
type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Customer, error)
	Update(ctx context.Context, customer *domain.Customer) error
	List(ctx context.Context, page, limit int) ([]domain.Customer, int64, error)
}

// AddressRepository handles customer address persistence.
type AddressRepository interface {
	Create(ctx context.Context, address *domain.Address) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error)
	ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Address, error)
	Update(ctx context.Context, address *domain.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetDefault(ctx context.Context, customerID, addressID uuid.UUID) error
}

// ProductRepository handles product persistence.
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter domain.ProductListFilter) ([]domain.Product, int64, error)
	DecrementStock(ctx context.Context, id uuid.UUID, quantity int) error
	IncrementStock(ctx context.Context, id uuid.UUID, quantity int) error
}

// ProductImageRepository handles product image persistence.
type ProductImageRepository interface {
	Create(ctx context.Context, image *domain.ProductImage) error
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]domain.ProductImage, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SetPrimary(ctx context.Context, productID, imageID uuid.UUID) error
}

// CategoryRepository handles category persistence.
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, activeOnly bool) ([]domain.Category, error)
}

// BannerRepository handles banner persistence.
type BannerRepository interface {
	Create(ctx context.Context, banner *domain.Banner) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Banner, error)
	Update(ctx context.Context, banner *domain.Banner) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	ListActive(ctx context.Context) ([]domain.Banner, error)
	List(ctx context.Context) ([]domain.Banner, error)
}

// PromoRepository handles promo code persistence.
type PromoRepository interface {
	Create(ctx context.Context, promo *domain.PromoCode) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PromoCode, error)
	GetByCode(ctx context.Context, code string) (*domain.PromoCode, error)
	Update(ctx context.Context, promo *domain.PromoCode) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	IncrementUsage(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]domain.PromoCode, error)
}

// CartRepository handles cart item persistence.
type CartRepository interface {
	AddItem(ctx context.Context, item *domain.CartItem) error
	GetItem(ctx context.Context, customerID, productID uuid.UUID) (*domain.CartItem, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*domain.CartItem, error)
	UpdateItem(ctx context.Context, item *domain.CartItem) error
	RemoveItem(ctx context.Context, id uuid.UUID) error
	ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.CartItem, error)
	ClearCart(ctx context.Context, customerID uuid.UUID) error
}

// OrderRepository handles order persistence.
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	CreateItems(ctx context.Context, items []domain.OrderItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
	Update(ctx context.Context, order *domain.Order) error
	List(ctx context.Context, filter domain.OrderListFilter) ([]domain.Order, int64, error)
	SalesReport(ctx context.Context, startDate, endDate time.Time) (*domain.SalesReport, error)
}

// PaymentRepository handles payment persistence.
type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error)
	GetByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
}
