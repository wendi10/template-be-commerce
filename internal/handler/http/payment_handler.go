package http

import (
	"encoding/json"
	"io"
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

type PaymentHandler struct {
	paymentUC usecase.PaymentUseCase
}

func NewPaymentHandler(paymentUC usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{paymentUC: paymentUC}
}

// CreatePayment handles POST /api/v1/payments
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	var req domain.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid JSON body")
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	result, err := h.paymentUC.CreatePayment(r.Context(), userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "payment created", result)
}

// GetPaymentByOrder handles GET /api/v1/payments/order/{orderID}
func (h *PaymentHandler) GetPaymentByOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := uuid.Parse(chi.URLParam(r, "orderID"))
	if err != nil {
		response.ValidationError(w, "invalid order ID")
		return
	}

	pay, err := h.paymentUC.GetPaymentByOrderID(r.Context(), orderID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "payment fetched", pay)
}

// HandleCallback handles POST /api/v1/payments/callback/{provider}
// This endpoint is called by the payment gateway webhook.
func (h *PaymentHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	providerStr := chi.URLParam(r, "provider")
	provider := domain.PaymentProvider(providerStr)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		response.ValidationError(w, "failed to read request body")
		return
	}

	// Collect headers for signature verification
	headers := make(map[string]string)
	for key := range r.Header {
		headers[key] = r.Header.Get(key)
	}

	if err := h.paymentUC.HandleCallback(r.Context(), provider, payload, headers); err != nil {
		// Return 200 even on error to prevent gateway retries for business logic errors
		// Log happens inside the use case
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
}
