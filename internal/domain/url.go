package domain

import "time"

// URL represents the URL entry in our system
type URL struct {
	ID          string    `json:"id" db:"id"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	VisitCount  string    `json:"visit_count" db:"visit_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	LastVisited time.Time `json:"last_visited" db:"last_visited"`
}

// CreateURLRequest represents the request body for creating a new URL
type CreateURLRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	CustomAlias string `json:"custom_alias,omitempty" validate:"omitempty,alphanum"`
}

// URLResponse represents the response for URL operations
type URLResponse struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	VisitCount  string `json:"visit_count"`
	CreatedAt   string `json:"created_at"`
	ShortURL    string `json:"short_url"`
}
