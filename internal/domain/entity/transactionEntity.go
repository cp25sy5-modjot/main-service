package entity

import (
	"time"
)

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

	// Relationships
	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category,omitempty"`
}

// Next Release: Split Transaction into Transaction and TransactionItem
// This will allow multiple items per transaction in the future
// Remove Quantity from TransactionItem as well
// type Transaction struct {
// 	TransactionID string    `gorm:"primaryKey;autoIncrement:false" json:"transaction_id" validate:"required"`
// 	UserID        string    `json:"user_id" validate:"required"`
// 	Date          time.Time `json:"date" validate:"required"`
// 	Type          string    `json:"type" validate:"required"` // e.g. "manual", "upload"

// 	// Optional summary fields
// 	TotalPrice  float64 `json:"total_price,omitempty"`
// 	// Relationships
// 	Items []TransactionItem `gorm:"foreignKey:TransactionID;references:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items,omitempty"`
// }

// type TransactionItem struct {
// 	TransactionID string `gorm:"primaryKey;autoIncrement:false" json:"transaction_id" validate:"required"`
// 	ItemID        string `gorm:"primaryKey;autoIncrement:false" json:"item_id" validate:"required"`

// 	Title      string  `json:"title" validate:"required"`
// 	Price      float64 `json:"price" validate:"required"`
// 	CategoryID string  `json:"category_id" validate:"required"`

// 	// Relationships
// 	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category,omitempty"`

// 	// Back-reference (optional but nice to have)
// 	Transaction Transaction `gorm:"foreignKey:TransactionID;references:TransactionID" json:"-"`
// }
