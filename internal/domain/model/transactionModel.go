package model

import (
	"time"
)

type TransactionSearchParams struct {
	TransactionID string `json:"transaction_id" validate:"required"`
	ItemID        string `json:"item_id" validate:"required"`
	UserID        string `json:"user_id" validate:"required"`
}

type TransactionInsertReq struct {
	Title      string    `json:"title" validate:"required,min=2,max=50"`
	Price      float64   `json:"price" validate:"required"`
	Quantity   float64   `json:"quantity" validate:"required"`
	CategoryID *string   `json:"category_id"`
	Date       time.Time `json:"date"`
}

type TransactionUpdateReq struct {
	Title      string    `json:"title" validate:"required,min=2,max=50"`
	Price      float64   `json:"price" validate:"required"`
	Quantity   float64   `json:"quantity" validate:"required"`
	CategoryID *string   `json:"category_id"`
	Date       time.Time `json:"date" validate:"required"`
}

type TransactionRes struct {
	TransactionID     string    `json:"transaction_id"`
	ItemID            string    `json:"item_id"`
	Title             string    `json:"title"`
	Price             float64   `json:"price"`
	Quantity          float64   `json:"quantity"`
	TotalPrice        float64   `json:"total_price"`
	Date              time.Time `json:"date"`
	Type              string    `json:"type"`
	CategoryID        *string   `json:"category_id"`
	CategoryName      string    `json:"category_name"`
	CategoryColorCode string    `json:"category_color_code"`
}

type TransactionFilter struct {
	Date *time.Time `json:"date"`
}
