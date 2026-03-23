package entity

import (
	"time"
)

type TransactionType string

const (
	TransactionManual  TransactionType = "manual"
	TransactionUpload  TransactionType = "upload"
	TransactionFixCost TransactionType = "fix_cost"
)

type Transaction struct {
	TransactionID string `gorm:"primaryKey;autoIncrement:false" `
	UserID        string
	Title         string
	Date          time.Time
	Type          TransactionType

	//fixcost
	RunDate   *time.Time `gorm:"type:date;index"`
	FixCostID *string

	// Relationships
	Items []TransactionItem `gorm:"foreignKey:TransactionID;references:TransactionID"`
}
