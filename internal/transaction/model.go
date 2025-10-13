package transaction

import (
	"time"
)

type TransactionRes struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Price     float64   `json:"price"`
	Date      time.Time `json:"date"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionReq struct {
	Title    string  `json:"title" validate:"required,min=2,max=50"`
	Price    float64 `json:"price" validate:"required"`
	Category string  `json:"category"`
}
