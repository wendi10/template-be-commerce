package payment

import (
	"context"

	"github.com/template-be-commerce/internal/domain"
)

// Gateway is the abstraction interface for all payment providers.
// Any provider (Doku, Midtrans, Xendit, etc.) must implement this interface.
type Gateway interface {
	// CreateTransaction initiates a payment transaction.
	CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*TransactionResponse, error)

	// HandleCallback processes the payment callback/webhook payload.
	HandleCallback(ctx context.Context, payload []byte, headers map[string]string) (*CallbackResult, error)

	// ProviderName returns the string identifier of the provider.
	ProviderName() domain.PaymentProvider
}

// CreateTransactionRequest is the standardised input for creating a payment.
type CreateTransactionRequest struct {
	OrderID       string
	OrderNumber   string
	Amount        string // string for decimal precision
	PaymentMethod string
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	Description   string
	CallbackURL   string
	ReturnURL     string
	ExpiredAt     int64 // Unix timestamp
}

// TransactionResponse is the normalised output after creating a payment transaction.
type TransactionResponse struct {
	TransactionID string
	PaymentURL    string
	ExpiredAt     int64
	RawResponse   string
}

// CallbackResult is the normalised result from parsing a provider callback.
type CallbackResult struct {
	TransactionID string
	OrderID       string
	Status        domain.PaymentStatus
	Amount        string
	RawPayload    string
}
