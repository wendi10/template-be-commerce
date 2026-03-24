package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/pkg/apperrors"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	// Create order without associations; items are created via CreateItems
	return r.db.WithContext(ctx).
		Omit("Items", "Customer", "Address", "PromoCode", "Payment").
		Create(order).Error
}

func (r *orderRepository) CreateItems(ctx context.Context, items []domain.OrderItem) error {
	for i := range items {
		if items[i].ID == uuid.Nil {
			items[i].ID = uuid.New()
		}
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Address").
		Preload("PromoCode").
		Where("id = ?", id).
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Address").
		Where("order_number = ?", orderNumber).
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.Order{}).
		Where("id = ?", id).
		Update("status", string(status)).Error
}

func (r *orderRepository) Update(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).
		Omit("Items", "Customer", "Address", "PromoCode", "Payment").
		Save(order).Error
}

func (r *orderRepository) List(ctx context.Context, filter domain.OrderListFilter) ([]domain.Order, int64, error) {
	db := r.db.WithContext(ctx).Model(&domain.Order{})

	if filter.CustomerID != nil {
		db = db.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.Status != nil {
		db = db.Where("status = ?", string(*filter.Status))
	}
	if filter.StartDate != nil {
		db = db.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		db = db.Where("created_at <= ?", *filter.EndDate)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []domain.Order
	offset := (filter.Page - 1) * filter.Limit
	err := db.
		Preload("Items").
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(offset).
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) SalesReport(ctx context.Context, startDate, endDate time.Time) (*domain.SalesReport, error) {
	type agg struct {
		TotalOrders     int64
		TotalRevenue    string
		AverageOrder    string
		PendingOrders   int64
		CompletedOrders int64
		CancelledOrders int64
	}

	var a agg
	err := r.db.WithContext(ctx).Model(&domain.Order{}).
		Select(`
			COUNT(*) AS total_orders,
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COALESCE(AVG(total_amount), 0) AS average_order,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) AS pending_orders,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) AS completed_orders,
			SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) AS cancelled_orders
		`).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&a).Error
	if err != nil {
		return nil, err
	}

	// Top-selling products
	type topRow struct {
		ProductID   string
		ProductName string
		TotalSold   int64
		Revenue     string
	}
	var rows []topRow
	_ = r.db.WithContext(ctx).Model(&domain.OrderItem{}).
		Select("product_id, product_name, SUM(quantity) AS total_sold, SUM(subtotal) AS revenue").
		Joins("JOIN orders o ON o.id = order_items.order_id").
		Where("o.created_at BETWEEN ? AND ? AND o.status NOT IN ('pending','cancelled')", startDate, endDate).
		Group("product_id, product_name").
		Order("total_sold DESC").
		Limit(10).
		Scan(&rows).Error

	var topProducts []domain.TopProduct
	for _, row := range rows {
		pid, _ := uuid.Parse(row.ProductID)
		tp := domain.TopProduct{
			ProductID:   pid,
			ProductName: row.ProductName,
			TotalSold:   row.TotalSold,
		}
		_ = tp.Revenue.Scan(row.Revenue)
		topProducts = append(topProducts, tp)
	}

	report := &domain.SalesReport{
		TotalOrders:     a.TotalOrders,
		PendingOrders:   a.PendingOrders,
		CompletedOrders: a.CompletedOrders,
		CancelledOrders: a.CancelledOrders,
		TopProducts:     topProducts,
		StartDate:       startDate,
		EndDate:         endDate,
	}
	_ = report.TotalRevenue.Scan(a.TotalRevenue)
	_ = report.AverageOrder.Scan(a.AverageOrder)

	return report, nil
}
