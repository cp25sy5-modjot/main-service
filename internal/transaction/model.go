package transaction

import (
	"time"
)

type SearchParams struct {
	TransactionID string `json:"transaction_id" validate:"required"`
	ProductID     string `json:"product_id" validate:"required"`
	UserID        string `json:"user_id" validate:"required"`
}

type TransactionInsertReq struct {
	Title    string  `json:"title" validate:"required,min=2,max=50"`
	Price    float64 `json:"price" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
	Category string  `json:"category" validate:"required"`
}

type TransactionUpdateReq struct {
	TransactionInsertReq
	Date string `json:"date" validate:"required"`
}

type TransactionRes struct {
	TransactionID string    `gorm:"primaryKey" json:"transaction_id"`
	ProductID     string    `json:"product_id"`
	UserID        string    `json:"user_id"`
	Title         string    `json:"title"`
	Price         float64   `json:"price"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`
	Category      string    `json:"category"`
	CreatedAt     time.Time `json:"created_at"`
}
