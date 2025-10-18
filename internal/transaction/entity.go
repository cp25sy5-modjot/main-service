package transaction

import "time"

type Transaction struct {
	TransactionID string    `gorm:"primaryKey;autoIncrement:false" json:"transaction_id" validate:"required"`
	ProductID     string    `gorm:"primaryKey;autoIncrement:false" json:"product_id" validate:"required"`
	UserID        string    `json:"user_id" validate:"required"`
	Title         string    `json:"title" validate:"required"`
	Price         float64   `json:"price" validate:"required"`
	Amount        float64   `json:"amount" validate:"required"`
	Date          time.Time `json:"date" validate:"required"`
	Type          string    `json:"type" validate:"required"`
	Category      string    `json:"category" validate:"required"`
}
