package transaction

import "time"

type Transaction struct {
	TransactionID string    `gorm:"primaryKey;autoIncrement:false" json:"transaction_id"`
	ProductID     string    `gorm:"primaryKey;autoIncrement:false" json:"product_id"`
	UserID        string    `json:"user_id"`
	Title         string    `json:"title"`
	Price         float64   `json:"price"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`
	Category      string    `json:"category"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
