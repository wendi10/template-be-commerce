package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/middleware"
	"github.com/template-be-commerce/internal/usecase"
	"github.com/template-be-commerce/pkg/apperrors"
	"github.com/template-be-commerce/pkg/pagination"
	"github.com/template-be-commerce/pkg/response"
	"github.com/template-be-commerce/pkg/validator"
)

type OrderHandler struct {
	orderUC usecase.OrderUseCase
}

func NewOrderHandler(orderUC usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUC: orderUC}
}

// CreateOrder handles POST /api/v1/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	var req domain.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	order, err := h.orderUC.CreateOrder(r.Context(), userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "order created", order)
}

// GetOrder handles GET /api/v1/orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	orderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid order ID")
		return
	}

	order, err := h.orderUC.GetOrder(r.Context(), userID, orderID, false)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "order fetched", order)
}

// ListMyOrders handles GET /api/v1/orders
func (h *OrderHandler) ListMyOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	p := pagination.FromRequest(r)
	filter := domain.OrderListFilter{Page: p.Page, Limit: p.Limit}

	if status := r.URL.Query().Get("status"); status != "" {
		s := domain.OrderStatus(status)
		filter.Status = &s
	}

	result, err := h.orderUC.ListOrders(r.Context(), userID, filter, false)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Paginated(w, "orders fetched", result.Orders, response.Meta{
		Page: result.Page, Limit: result.Limit, Total: result.Total, TotalPages: result.TotalPages,
	})
}

// CancelOrder handles POST /api/v1/orders/{id}/cancel
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	orderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid order ID")
		return
	}

	if err := h.orderUC.CancelOrder(r.Context(), userID, orderID); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "order cancelled", nil)
}

// --- Admin Handlers ---

// AdminListOrders handles GET /api/v1/admin/orders
func (h *OrderHandler) AdminListOrders(w http.ResponseWriter, r *http.Request) {
	p := pagination.FromRequest(r)
	filter := domain.OrderListFilter{Page: p.Page, Limit: p.Limit}

	q := r.URL.Query()
	if status := q.Get("status"); status != "" {
		s := domain.OrderStatus(status)
		filter.Status = &s
	}
	if startStr := q.Get("start_date"); startStr != "" {
		if t, err := time.Parse(time.DateOnly, startStr); err == nil {
			// Start of day
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			filter.StartDate = &t
		}
	}
	if endStr := q.Get("end_date"); endStr != "" {
		if t, err := time.Parse(time.DateOnly, endStr); err == nil {
			// End of day — include all orders on this calendar date
			t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
			filter.EndDate = &t
		}
	}

	userID, _ := middleware.UserIDFromContext(r.Context())
	result, err := h.orderUC.ListOrders(r.Context(), userID, filter, true)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Paginated(w, "orders fetched", result.Orders, response.Meta{
		Page: result.Page, Limit: result.Limit, Total: result.Total, TotalPages: result.TotalPages,
	})
}

// AdminGetOrder handles GET /api/v1/admin/orders/{id}
func (h *OrderHandler) AdminGetOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid order ID")
		return
	}

	userID, _ := middleware.UserIDFromContext(r.Context())
	order, err := h.orderUC.GetOrder(r.Context(), userID, orderID, true)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "order fetched", order)
}

// AdminUpdateOrderStatus handles PATCH /api/v1/admin/orders/{id}/status
func (h *OrderHandler) AdminUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid order ID")
		return
	}

	var req domain.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	order, err := h.orderUC.UpdateOrderStatus(r.Context(), orderID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "order status updated", order)
}

// AdminSalesReport handles GET /api/v1/admin/reports/sales
func (h *OrderHandler) AdminSalesReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day()-30, 0, 0, 0, 0, now.Location())
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	if s := q.Get("start_date"); s != "" {
		if t, err := time.Parse(time.DateOnly, s); err == nil {
			// Start of day
			startDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		}
	}
	if e := q.Get("end_date"); e != "" {
		if t, err := time.Parse(time.DateOnly, e); err == nil {
			// End of day — include all orders created on this calendar date
			endDate = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
		}
	}

	report, err := h.orderUC.GetSalesReport(r.Context(), startDate, endDate)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "sales report fetched", report)
}
