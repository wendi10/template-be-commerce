package pagination

import (
	"math"
	"net/http"
	"strconv"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

type Params struct {
	Page  int
	Limit int
}

// FromRequest extracts pagination parameters from query string.
func FromRequest(r *http.Request) Params {
	page := DefaultPage
	limit := DefaultLimit

	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			if v > MaxLimit {
				v = MaxLimit
			}
			limit = v
		}
	}

	return Params{Page: page, Limit: limit}
}

// Offset calculates the SQL OFFSET value.
func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

// TotalPages calculates the number of pages.
func TotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}
