package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/template-be-commerce/pkg/apperrors"
	"github.com/template-be-commerce/pkg/logger"
	"go.uber.org/zap"
)

// Envelope wraps all API responses for consistency.
type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrDetail  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Meta holds pagination info.
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("failed to encode response", zap.Error(err))
	}
}

func Success(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, Envelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusCreated, Envelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Paginated(w http.ResponseWriter, message string, data interface{}, meta Meta) {
	JSON(w, http.StatusOK, Envelope{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

func Error(w http.ResponseWriter, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		JSON(w, appErr.Code, Envelope{
			Success: false,
			Error: &ErrDetail{
				Code:    appErr.Code,
				Message: appErr.Message,
				Detail:  appErr.Detail,
			},
		})
		return
	}

	// Unhandled error → 500
	logger.Error("unhandled error", zap.Error(err))
	JSON(w, http.StatusInternalServerError, Envelope{
		Success: false,
		Error: &ErrDetail{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		},
	})
}

func ValidationError(w http.ResponseWriter, detail string) {
	JSON(w, http.StatusBadRequest, Envelope{
		Success: false,
		Error: &ErrDetail{
			Code:    http.StatusBadRequest,
			Message: "validation error",
			Detail:  detail,
		},
	})
}
