package receipt

import "time"

type ReceiptRes struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}
