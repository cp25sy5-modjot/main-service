package entity

import (
	"time"
)

type Transaction struct {
	TransactionID string `gorm:"primaryKey;autoIncrement:false"`
	ItemID        string `gorm:"primaryKey;autoIncrement:false"`
	UserID        string
	Title         string `gorm:"length:20"`
	Price         float64
	Quantity      float64
	Date          time.Time
	Type          string

	// ทำ nullable ไว้ ถ้าอยากใช้ OnDelete:SET NULL เวลา category ถูกลบ
	CategoryID *string

	// belongs-to User (เพื่อให้ DB มี FK user_id ด้วย)
	User User `gorm:"foreignKey:UserID;references:UserID"`

	// belongs-to Category (composite key)
	Category Category `gorm:"foreignKey:UserID,CategoryID;references:UserID,CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Next Release: Split Transaction into Transaction and TransactionItem
// This will allow multiple items per transaction in the future
// Remove Quantity from TransactionItem as well
// type Transaction struct {
// 	TransactionID string    `gorm:"primaryKey;autoIncrement:false" `
// 	UserID        string
// 	Date          time.Time
// 	Type          string    // e.g. "manual", "upload"

// 	// Optional summary fields
// 	TotalPrice  float64
// 	// Relationships
// 	Items []TransactionItem `gorm:"foreignKey:TransactionID;references:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
// }

// type TransactionItem struct {
// 	TransactionID string `gorm:"primaryKey;autoIncrement:false" `
// 	ItemID        string `gorm:"primaryKey;autoIncrement:false" `

// 	Title      string
// 	Price      float64
// 	CategoryID string

//	// Relationships
// 	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID"`

// 	// Back-reference (optional but nice to have)
// 	Transaction Transaction `gorm:"foreignKey:TransactionID;references:TransactionID"`
// }
