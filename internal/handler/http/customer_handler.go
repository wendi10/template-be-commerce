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

type CustomerHandler struct {
	customerUC usecase.CustomerUseCase
}

func NewCustomerHandler(customerUC usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{customerUC: customerUC}
}

// GetProfile handles GET /api/v1/me
func (h *CustomerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	customer, err := h.customerUC.GetProfile(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "profile fetched", customer)
}

// UpdateProfile handles PUT /api/v1/me
func (h *CustomerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	var req domain.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	customer, err := h.customerUC.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "profile updated", customer)
}

// ListAddresses handles GET /api/v1/me/addresses
func (h *CustomerHandler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	addresses, err := h.customerUC.ListAddresses(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "addresses fetched", addresses)
}

// CreateAddress handles POST /api/v1/me/addresses
func (h *CustomerHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	var req domain.CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	address, err := h.customerUC.CreateAddress(r.Context(), userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "address created", address)
}

// UpdateAddress handles PUT /api/v1/me/addresses/{addressID}
func (h *CustomerHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	addressID, err := uuid.Parse(chi.URLParam(r, "addressID"))
	if err != nil {
		response.ValidationError(w, "invalid address ID")
		return
	}

	var req domain.UpdateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}

	address, err := h.customerUC.UpdateAddress(r.Context(), userID, addressID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "address updated", address)
}

// DeleteAddress handles DELETE /api/v1/me/addresses/{addressID}
func (h *CustomerHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	addressID, err := uuid.Parse(chi.URLParam(r, "addressID"))
	if err != nil {
		response.ValidationError(w, "invalid address ID")
		return
	}

	if err := h.customerUC.DeleteAddress(r.Context(), userID, addressID); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}

// SetDefaultAddress handles PATCH /api/v1/me/addresses/{addressID}/default
func (h *CustomerHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	addressID, err := uuid.Parse(chi.URLParam(r, "addressID"))
	if err != nil {
		response.ValidationError(w, "invalid address ID")
		return
	}

	if err := h.customerUC.SetDefaultAddress(r.Context(), userID, addressID); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "default address set", nil)
}
