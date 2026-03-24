package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/usecase"
	"github.com/template-be-commerce/pkg/response"
	"github.com/template-be-commerce/pkg/validator"
)

type BannerHandler struct {
	bannerUC usecase.BannerUseCase
}

func NewBannerHandler(bannerUC usecase.BannerUseCase) *BannerHandler {
	return &BannerHandler{bannerUC: bannerUC}
}

// ListActiveBanners handles GET /api/v1/banners
func (h *BannerHandler) ListActiveBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := h.bannerUC.ListActive(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "banners fetched", banners)
}

// AdminListBanners handles GET /api/v1/admin/banners
func (h *BannerHandler) AdminListBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := h.bannerUC.ListAll(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "banners fetched", banners)
}

// AdminCreateBanner handles POST /api/v1/admin/banners
func (h *BannerHandler) AdminCreateBanner(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	banner, err := h.bannerUC.Create(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "banner created", banner)
}

// AdminGetBanner handles GET /api/v1/admin/banners/{id}
func (h *BannerHandler) AdminGetBanner(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid banner ID")
		return
	}

	banner, err := h.bannerUC.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "banner fetched", banner)
}

// AdminUpdateBanner handles PUT /api/v1/admin/banners/{id}
func (h *BannerHandler) AdminUpdateBanner(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid banner ID")
		return
	}

	var req domain.UpdateBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}

	banner, err := h.bannerUC.Update(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "banner updated", banner)
}

// AdminDeleteBanner handles DELETE /api/v1/admin/banners/{id}
func (h *BannerHandler) AdminDeleteBanner(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid banner ID")
		return
	}

	if err := h.bannerUC.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
