package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/usecase"
	"github.com/template-be-commerce/pkg/response"
	"github.com/template-be-commerce/pkg/validator"
)

type PromoHandler struct {
	promoUC usecase.PromoUseCase
}

func NewPromoHandler(promoUC usecase.PromoUseCase) *PromoHandler {
	return &PromoHandler{promoUC: promoUC}
}

// ValidatePromo handles POST /api/v1/promos/validate
func (h *PromoHandler) ValidatePromo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code     string          `json:"code" validate:"required"`
		SubTotal decimal.Decimal `json:"sub_total" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	result, _ := h.promoUC.ValidateAndCalculate(r.Context(), req.Code, req.SubTotal)
	response.Success(w, "promo validated", result)
}

// AdminListPromos handles GET /api/v1/admin/promos
func (h *PromoHandler) AdminListPromos(w http.ResponseWriter, r *http.Request) {
	promos, err := h.promoUC.List(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "promos fetched", promos)
}

// AdminGetPromo handles GET /api/v1/admin/promos/{id}
func (h *PromoHandler) AdminGetPromo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid promo ID")
		return
	}

	promo, err := h.promoUC.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "promo fetched", promo)
}

// AdminCreatePromo handles POST /api/v1/admin/promos
func (h *PromoHandler) AdminCreatePromo(w http.ResponseWriter, r *http.Request) {
	var req domain.CreatePromoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	promo, err := h.promoUC.Create(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "promo created", promo)
}

// AdminUpdatePromo handles PUT /api/v1/admin/promos/{id}
func (h *PromoHandler) AdminUpdatePromo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid promo ID")
		return
	}

	var req domain.UpdatePromoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}

	promo, err := h.promoUC.Update(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "promo updated", promo)
}

// AdminDeletePromo handles DELETE /api/v1/admin/promos/{id}
func (h *PromoHandler) AdminDeletePromo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid promo ID")
		return
	}

	if err := h.promoUC.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
