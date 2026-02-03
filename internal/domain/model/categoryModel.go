package model

import "time"

type CategoryReq struct {
	CategoryName string  `json:"category_name" validate:"required,min=1,max=50"`
	Budget       float64 `json:"budget" validate:"required,min=0"`
	ColorCode    string  `json:"color_code" validate:"required,min=7,max=7"`
	Icon         string  `json:"icon" validate:"required"`
}

type CategoryUpdateReq struct {
	CategoryName string  `json:"category_name" validate:"required,min=1,max=50"`
	Budget       float64 `json:"budget" validate:"required,min=0"`
	ColorCode    string  `json:"color_code" validate:"required,min=7,max=7"`
	Icon         string  `json:"icon" validate:"required"`
}

type CategoryRes struct {
	CategoryID   *string   `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Budget       float64   `json:"budget"`
	ColorCode    string    `json:"color_code"`
	Icon         string    `json:"icon"`
	CreatedAt    time.Time `json:"created_at"`

	BudgetUsage float64 `json:"budget_usage"`
}

type CategorySearchParams struct {
	CategoryID string
	UserID     string
}
