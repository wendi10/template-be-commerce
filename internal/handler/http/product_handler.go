package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/usecase"
	"github.com/template-be-commerce/pkg/pagination"
	"github.com/template-be-commerce/pkg/response"
	"github.com/template-be-commerce/pkg/validator"
)

type ProductHandler struct {
	productUC usecase.ProductUseCase
}

func NewProductHandler(productUC usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUC: productUC}
}

// --- Public ---

// ListProducts handles GET /api/v1/products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	p := pagination.FromRequest(r)
	q := r.URL.Query()

	isActiveTrue := true
	filter := domain.ProductListFilter{
		Search:    q.Get("search"),
		Page:      p.Page,
		Limit:     p.Limit,
		SortBy:    q.Get("sort_by"),
		SortOrder: q.Get("sort_order"),
		IsActive:  &isActiveTrue,
	}

	if catStr := q.Get("category_id"); catStr != "" {
		if catID, err := uuid.Parse(catStr); err == nil {
			filter.CategoryID = &catID
		}
	}
	if minStr := q.Get("min_price"); minStr != "" {
		if v, err := decimal.NewFromString(minStr); err == nil {
			filter.MinPrice = &v
		}
	}
	if maxStr := q.Get("max_price"); maxStr != "" {
		if v, err := decimal.NewFromString(maxStr); err == nil {
			filter.MaxPrice = &v
		}
	}

	result, err := h.productUC.ListProducts(r.Context(), filter)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Paginated(w, "products fetched", result.Products, response.Meta{
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: result.TotalPages,
	})
}

// GetProduct handles GET /api/v1/products/{id}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid product ID")
		return
	}

	product, err := h.productUC.GetProduct(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "product fetched", product)
}

// GetProductBySlug handles GET /api/v1/products/slug/{slug}
func (h *ProductHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	s := chi.URLParam(r, "slug")
	product, err := h.productUC.GetProductBySlug(r.Context(), s)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "product fetched", product)
}

// ListCategories handles GET /api/v1/categories
func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.productUC.ListCategories(r.Context(), true)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "categories fetched", categories)
}

// --- Admin ---

// AdminListProducts handles GET /api/v1/admin/products
func (h *ProductHandler) AdminListProducts(w http.ResponseWriter, r *http.Request) {
	p := pagination.FromRequest(r)
	q := r.URL.Query()

	filter := domain.ProductListFilter{
		Search:    q.Get("search"),
		Page:      p.Page,
		Limit:     p.Limit,
		SortBy:    q.Get("sort_by"),
		SortOrder: q.Get("sort_order"),
	}

	result, err := h.productUC.ListProducts(r.Context(), filter)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Paginated(w, "products fetched", result.Products, response.Meta{
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: result.TotalPages,
	})
}

// AdminCreateProduct handles POST /api/v1/admin/products
func (h *ProductHandler) AdminCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	product, err := h.productUC.CreateProduct(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "product created", product)
}

// AdminUpdateProduct handles PUT /api/v1/admin/products/{id}
func (h *ProductHandler) AdminUpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid product ID")
		return
	}

	var req domain.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}

	product, err := h.productUC.UpdateProduct(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "product updated", product)
}

// AdminDeleteProduct handles DELETE /api/v1/admin/products/{id}
func (h *ProductHandler) AdminDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid product ID")
		return
	}

	if err := h.productUC.DeleteProduct(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}

// AdminAddProductImage handles POST /api/v1/admin/products/{id}/images
func (h *ProductHandler) AdminAddProductImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid product ID")
		return
	}

	var req domain.AddProductImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	image, err := h.productUC.AddProductImage(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "image added", image)
}

// AdminDeleteProductImage handles DELETE /api/v1/admin/products/{id}/images/{imageID}
func (h *ProductHandler) AdminDeleteProductImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid product ID")
		return
	}
	imageID, err := uuid.Parse(chi.URLParam(r, "imageID"))
	if err != nil {
		response.ValidationError(w, "invalid image ID")
		return
	}

	if err := h.productUC.DeleteProductImage(r.Context(), id, imageID); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}

// --- Category Admin ---

// AdminCreateCategory handles POST /api/v1/admin/categories
func (h *ProductHandler) AdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	category, err := h.productUC.CreateCategory(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "category created", category)
}

// AdminUpdateCategory handles PUT /api/v1/admin/categories/{id}
func (h *ProductHandler) AdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid category ID")
		return
	}

	var req domain.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}

	category, err := h.productUC.UpdateCategory(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "category updated", category)
}

// AdminDeleteCategory handles DELETE /api/v1/admin/categories/{id}
func (h *ProductHandler) AdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid category ID")
		return
	}

	if err := h.productUC.DeleteCategory(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}

// AdminListCategories handles GET /api/v1/admin/categories
func (h *ProductHandler) AdminListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.productUC.ListCategories(r.Context(), false)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "categories fetched", categories)
}
