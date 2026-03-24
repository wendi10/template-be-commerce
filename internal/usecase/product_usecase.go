package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/repository"
	"github.com/template-be-commerce/pkg/slug"
)

type ProductUseCase interface {
	// Admin
	CreateProduct(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req domain.UpdateProductRequest) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	AddProductImage(ctx context.Context, productID uuid.UUID, req domain.AddProductImageRequest) (*domain.ProductImage, error)
	DeleteProductImage(ctx context.Context, productID, imageID uuid.UUID) error
	SetPrimaryImage(ctx context.Context, productID, imageID uuid.UUID) error

	// Admin & Public
	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetProductBySlug(ctx context.Context, s string) (*domain.Product, error)
	ListProducts(ctx context.Context, filter domain.ProductListFilter) (*domain.ProductListResponse, error)

	// Categories
	CreateCategory(ctx context.Context, req domain.CreateCategoryRequest) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req domain.UpdateCategoryRequest) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context, activeOnly bool) ([]domain.Category, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error)
}

type productUseCase struct {
	productRepo      repository.ProductRepository
	productImageRepo repository.ProductImageRepository
	categoryRepo     repository.CategoryRepository
}

func NewProductUseCase(
	productRepo repository.ProductRepository,
	productImageRepo repository.ProductImageRepository,
	categoryRepo repository.CategoryRepository,
) ProductUseCase {
	return &productUseCase{
		productRepo:      productRepo,
		productImageRepo: productImageRepo,
		categoryRepo:     categoryRepo,
	}
}

func (uc *productUseCase) CreateProduct(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error) {
	s := req.Slug
	if s == "" {
		s = slug.Generate(req.Name)
	}

	product := &domain.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        s,
		Description: req.Description,
		Price:       req.Price,
		Weight:      req.Weight,
		Stock:       req.Stock,
		IsActive:    req.IsActive,
	}

	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (uc *productUseCase) UpdateProduct(ctx context.Context, id uuid.UUID, req domain.UpdateProductRequest) (*domain.Product, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Slug != "" {
		product.Slug = req.Slug
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Weight != nil {
		product.Weight = *req.Weight
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (uc *productUseCase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.productRepo.SoftDelete(ctx, id)
}

func (uc *productUseCase) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return uc.productRepo.GetByID(ctx, id)
}

func (uc *productUseCase) GetProductBySlug(ctx context.Context, s string) (*domain.Product, error) {
	return uc.productRepo.GetBySlug(ctx, s)
}

func (uc *productUseCase) ListProducts(ctx context.Context, filter domain.ProductListFilter) (*domain.ProductListResponse, error) {
	products, total, err := uc.productRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if filter.Limit > 0 {
		totalPages = int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	}

	return &domain.ProductListResponse{
		Products:   products,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (uc *productUseCase) AddProductImage(ctx context.Context, productID uuid.UUID, req domain.AddProductImageRequest) (*domain.ProductImage, error) {
	_, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	image := &domain.ProductImage{
		ProductID: productID,
		URL:       req.URL,
		AltText:   req.AltText,
		IsPrimary: req.IsPrimary,
		SortOrder: req.SortOrder,
	}

	if err := uc.productImageRepo.Create(ctx, image); err != nil {
		return nil, err
	}

	if req.IsPrimary {
		_ = uc.productImageRepo.SetPrimary(ctx, productID, image.ID)
	}

	return image, nil
}

func (uc *productUseCase) DeleteProductImage(ctx context.Context, productID, imageID uuid.UUID) error {
	return uc.productImageRepo.Delete(ctx, imageID)
}

func (uc *productUseCase) SetPrimaryImage(ctx context.Context, productID, imageID uuid.UUID) error {
	return uc.productImageRepo.SetPrimary(ctx, productID, imageID)
}

// --- Categories ---

func (uc *productUseCase) CreateCategory(ctx context.Context, req domain.CreateCategoryRequest) (*domain.Category, error) {
	s := req.Slug
	if s == "" {
		s = slug.Generate(req.Name)
	}

	category := &domain.Category{
		Name:        req.Name,
		Slug:        s,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (uc *productUseCase) UpdateCategory(ctx context.Context, id uuid.UUID, req domain.UpdateCategoryRequest) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (uc *productUseCase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.categoryRepo.SoftDelete(ctx, id)
}

func (uc *productUseCase) ListCategories(ctx context.Context, activeOnly bool) ([]domain.Category, error) {
	return uc.categoryRepo.List(ctx, activeOnly)
}

func (uc *productUseCase) GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return uc.categoryRepo.GetByID(ctx, id)
}
