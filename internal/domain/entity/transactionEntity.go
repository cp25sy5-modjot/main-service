package entity

import (
	"time"
)

type TransactionType string

const (
	TransactionManual TransactionType = "manual"
	TransactionUpload TransactionType = "upload"
)

type Transaction struct {
	TransactionID string `gorm:"primaryKey;autoIncrement:false" `
	UserID        string
	Date          time.Time
	Type          TransactionType

	// Relationships
	Items []TransactionItem `gorm:"foreignKey:TransactionID;references:TransactionID"`
}

