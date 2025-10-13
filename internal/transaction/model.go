package transaction

import (
	"time"
)

type TransactionInsertReq struct {
	Title    string  `json:"title" validate:"required,min=2,max=50"`
	Price    float64 `json:"price" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
	Type     string  `json:"type" validate:"required,oneof=manual upload"`
	Category string  `json:"category"`
}

type TransactionUpdateReq struct {
	Title    string  `json:"title" validate:"min=2,max=50"`
	Price    float64 `json:"price" validate:"omitempty"`
	Amount   float64 `json:"amount" validate:"omitempty"`
	Type     string  `json:"type" validate:"omitempty,oneof=manual upload"`
	Category string  `json:"category"`
	Date     string  `json:"date"`
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
