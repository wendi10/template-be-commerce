package mysql

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/pkg/apperrors"
)

// ─── Product Repository ───────────────────────────────────────────────────────

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var p domain.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("id = ?", id).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	var p domain.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("slug = ?", slug).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).
		Model(product).
		Omit("Images", "Category").
		Save(product).Error
}

func (r *productRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *productRepository) List(ctx context.Context, filter domain.ProductListFilter) ([]domain.Product, int64, error) {
	db := r.db.WithContext(ctx).Model(&domain.Product{}).Where("deleted_at IS NULL")

	if filter.CategoryID != nil {
		db = db.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	if filter.MinPrice != nil {
		db = db.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		db = db.Where("price <= ?", *filter.MaxPrice)
	}
	if filter.IsActive != nil {
		db = db.Where("is_active = ?", *filter.IsActive)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Allowlist sort columns
	sortCol := "created_at"
	allowed := map[string]bool{"created_at": true, "price": true, "name": true, "stock": true}
	if allowed[filter.SortBy] {
		sortCol = filter.SortBy
	}
	sortDir := "DESC"
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortDir = "ASC"
	}

	var products []domain.Product
	offset := (filter.Page - 1) * filter.Limit
	err := db.
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Order(sortCol + " " + sortDir).
		Limit(filter.Limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) DecrementStock(ctx context.Context, id uuid.UUID, quantity int) error {
	result := r.db.WithContext(ctx).Model(&domain.Product{}).
		Where("id = ? AND stock >= ? AND deleted_at IS NULL", id, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrInsufficientStock
	}
	return nil
}

func (r *productRepository) IncrementStock(ctx context.Context, id uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Model(&domain.Product{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

// ─── Product Image Repository ─────────────────────────────────────────────────

type productImageRepository struct {
	db *gorm.DB
}

func NewProductImageRepository(db *gorm.DB) *productImageRepository {
	return &productImageRepository{db: db}
}

func (r *productImageRepository) Create(ctx context.Context, image *domain.ProductImage) error {
	if image.ID == uuid.Nil {
		image.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *productImageRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]domain.ProductImage, error) {
	var images []domain.ProductImage
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sort_order ASC").
		Find(&images).Error
	return images, err
}

func (r *productImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.ProductImage{}).Error
}

func (r *productImageRepository) SetPrimary(ctx context.Context, productID, imageID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.ProductImage{}).
			Where("product_id = ?", productID).
			Update("is_primary", false).Error; err != nil {
			return err
		}
		return tx.Model(&domain.ProductImage{}).
			Where("id = ?", imageID).
			Update("is_primary", true).Error
	})
}

// ─── Category Repository ──────────────────────────────────────────────────────

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var c domain.Category
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var c domain.Category
	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).
		Omit(clause.Associations).
		Save(category).Error
}

func (r *categoryRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.Category{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *categoryRepository) List(ctx context.Context, activeOnly bool) ([]domain.Category, error) {
	var categories []domain.Category
	db := r.db.WithContext(ctx).Where("deleted_at IS NULL")
	if activeOnly {
		db = db.Where("is_active = ?", true)
	}
	err := db.Order("name ASC").Find(&categories).Error
	return categories, err
}
