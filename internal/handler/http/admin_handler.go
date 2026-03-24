package http

import (
	"net/http"

	"github.com/template-be-commerce/internal/usecase"
	"github.com/template-be-commerce/pkg/pagination"
	"github.com/template-be-commerce/pkg/response"
)

// AdminHandler consolidates admin-specific operations that don't fit in
// other domain handlers (e.g., customer management, dashboard stats).
type AdminHandler struct {
	customerUC usecase.CustomerUseCase
}

func NewAdminHandler(customerUC usecase.CustomerUseCase) *AdminHandler {
	return &AdminHandler{customerUC: customerUC}
}

// ListCustomers handles GET /api/v1/admin/customers
func (h *AdminHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	p := pagination.FromRequest(r)

	customers, total, err := h.customerUC.ListAll(r.Context(), p.Page, p.Limit)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Paginated(w, "customers fetched", customers, response.Meta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: pagination.TotalPages(total, p.Limit),
	})
}
