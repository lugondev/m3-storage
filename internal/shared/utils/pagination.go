package utils

import "math"

// PaginationQuery represents the query parameters for pagination.
type PaginationQuery struct {
	Page     int `query:"page" json:"page"`
	PageSize int `query:"page_size" json:"page_size"`
}

// Pagination represents pagination metadata.
type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasPrevious bool  `json:"has_previous"`
	HasNext     bool  `json:"has_next"`
}

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// ValidateAndSetDefaults validates and sets default values for pagination parameters.
func (q *PaginationQuery) ValidateAndSetDefaults() {
	if q.Page < 1 {
		q.Page = DefaultPage
	}

	if q.PageSize < 1 {
		q.PageSize = DefaultPageSize
	} else if q.PageSize > MaxPageSize {
		q.PageSize = MaxPageSize
	}
}

// GetOffset calculates the offset for database queries.
func (q *PaginationQuery) GetOffset() int {
	return (q.Page - 1) * q.PageSize
}

// GetLimit returns the limit for database queries.
func (q *PaginationQuery) GetLimit() int {
	return q.PageSize
}

// NewPagination creates a new Pagination instance.
func NewPagination(query PaginationQuery, totalItems int64) Pagination {
	totalPages := int(math.Ceil(float64(totalItems) / float64(query.PageSize)))

	return Pagination{
		CurrentPage: query.Page,
		PageSize:    query.PageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasPrevious: query.Page > 1,
		HasNext:     query.Page < totalPages,
	}
}
