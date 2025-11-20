package transaction

import (
	"time"
)

type TransactionSearchParams struct {
	TransactionID string `json:"transaction_id" validate:"required"`
	ItemID        string `json:"product_id" validate:"required"`
	UserID        string `json:"user_id" validate:"required"`
}

type TransactionInsertReq struct {
	Title      string  `json:"title" validate:"required,min=2,max=50"`
	Price      float64 `json:"price" validate:"required"`
	Quantity   float64 `json:"quantity" validate:"required"`
	CategoryId string  `json:"category_id" validate:"required"`
}

type TransactionUpdateReq struct {
	TransactionInsertReq
	Date string `json:"date" validate:"required"`
}

type TransactionRes struct {
	TransactionID string    `json:"transaction_id"`
	ItemID        string    `json:"product_id"`
	Title         string    `json:"title"`
	Price         float64   `json:"price"`
	Quantity      float64   `json:"quantity"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`
	CategoryID    string    `json:"category_id"`
}

type TransactionFilter struct {
	Date *time.Time `json:"date"`
}
