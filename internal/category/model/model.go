package category

import "time"

type CategoryReq struct {
	CategoryName string  `json:"category_name" validate:"required"`
	Budget       float64 `json:"budget" validate:"required"`
}

type CategoryRes struct {
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Budget       float64   `json:"budget"`
	CreatedAt    time.Time `json:"created_at"`
}

type CategoryUpdateReq struct {
	CategoryName string  `json:"category_name"`
	Budget       float64 `json:"budget"`
}

type CategorySearchParams struct {
	CategoryID string
	UserID     string
}
