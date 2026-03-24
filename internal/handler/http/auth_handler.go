package http

import (
	"encoding/json"
	"net/http"

	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/usecase"
	"github.com/template-be-commerce/pkg/response"
	"github.com/template-be-commerce/pkg/validator"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	result, err := h.authUC.RegisterCustomer(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "registration successful", result)
}

// LoginCustomer handles POST /api/v1/auth/login
func (h *AuthHandler) LoginCustomer(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	result, err := h.authUC.LoginCustomer(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "login successful", result)
}

// RegisterAdmin handles POST /api/v1/admin/auth/register
func (h *AuthHandler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	result, err := h.authUC.RegisterAdmin(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "admin registration successful", result)
}

// LoginAdmin handles POST /api/v1/admin/auth/login
func (h *AuthHandler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	result, err := h.authUC.LoginAdmin(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "admin login successful", result)
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	result, err := h.authUC.RefreshToken(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "token refreshed", result)
}
