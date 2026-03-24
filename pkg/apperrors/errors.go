package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError is a domain-level error that carries an HTTP status code.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Constructors

func New(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func BadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func Forbidden(msg string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg}
}

func UnprocessableEntity(msg string) *AppError {
	return &AppError{Code: http.StatusUnprocessableEntity, Message: msg}
}

func InternalServer(msg string, err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg, Err: err}
}

// Is checks if the target error has a matching code.
func Is(err error, code int) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// Sentinel errors
var (
	ErrNotFound         = NotFound("resource not found")
	ErrUnauthorized     = Unauthorized("unauthorized")
	ErrForbidden        = Forbidden("forbidden")
	ErrInvalidToken     = Unauthorized("invalid or expired token")
	ErrEmailExists      = Conflict("email already registered")
	ErrInvalidPassword  = Unauthorized("invalid email or password")
	ErrInactiveAccount  = Unauthorized("account is inactive")
	ErrInsufficientStock = BadRequest("insufficient product stock")
	ErrCartEmpty        = BadRequest("cart is empty")
	ErrPromoExpired     = BadRequest("promo code is expired")
	ErrPromoUsageLimitReached = BadRequest("promo code usage limit reached")
	ErrInvalidPromo     = BadRequest("invalid promo code")
	ErrInvalidOrderTransition = BadRequest("invalid order status transition")
)
