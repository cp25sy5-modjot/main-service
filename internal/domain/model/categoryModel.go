package model

import "time"

type CategoryReq struct {
	CategoryName string  `json:"category_name" validate:"required"`
	Budget       float64 `json:"budget" validate:"required"`
	ColorCode    string  `json:"color_code" validate:"required"`
}

type CategoryRes struct {
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Budget       float64   `json:"budget"`
	ColorCode    string    `json:"color_code"`
	CreatedAt    time.Time `json:"created_at"`

	BudgetUsage float64 `json:"budget_usage,omitempty"`
}

type CategoryUpdateReq struct {
	CategoryName string  `json:"category_name"`
	Budget       float64 `json:"budget"`
	ColorCode    string  `json:"color_code"`
}

type CategorySearchParams struct {
	CategoryID string
	UserID     string
}
