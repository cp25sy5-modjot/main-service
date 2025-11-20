package transaction

import "time"

type Transaction struct {
	TransactionID string    `gorm:"primaryKey;autoIncrement:false" json:"transaction_id" validate:"required"`
	ItemID        string    `gorm:"primaryKey;autoIncrement:false" json:"item_id" validate:"required"`
	UserID        string    `json:"user_id" validate:"required"`
	Title         string    `json:"title" validate:"required"`
	Price         float64   `json:"price" validate:"required"`
	Quantity      float64   `json:"quantity" validate:"required"`
	Date          time.Time `json:"date" validate:"required"`
	Type          string    `json:"type" validate:"required"`
	CategoryID    string    `json:"category_id" validate:"required"`
}
