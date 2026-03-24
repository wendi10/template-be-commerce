package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/middleware"
	"github.com/template-be-commerce/internal/usecase"
	"github.com/template-be-commerce/pkg/apperrors"
	"github.com/template-be-commerce/pkg/response"
	"github.com/template-be-commerce/pkg/validator"
)

type CartHandler struct {
	cartUC usecase.CartUseCase
}

func NewCartHandler(cartUC usecase.CartUseCase) *CartHandler {
	return &CartHandler{cartUC: cartUC}
}

// GetCart handles GET /api/v1/cart
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	promoCode := r.URL.Query().Get("promo_code")
	summary, err := h.cartUC.GetCartSummary(r.Context(), userID, promoCode)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "cart fetched", summary)
}

// AddToCart handles POST /api/v1/cart
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	var req domain.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	item, err := h.cartUC.AddToCart(r.Context(), userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "item added to cart", item)
}

// UpdateCartItem handles PUT /api/v1/cart/{itemID}
func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		response.ValidationError(w, "invalid cart item ID")
		return
	}

	var req domain.UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	item, err := h.cartUC.UpdateCartItem(r.Context(), userID, itemID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "cart item updated", item)
}

// RemoveCartItem handles DELETE /api/v1/cart/{itemID}
func (h *CartHandler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		response.ValidationError(w, "invalid cart item ID")
		return
	}

	if err := h.cartUC.RemoveCartItem(r.Context(), userID, itemID); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
