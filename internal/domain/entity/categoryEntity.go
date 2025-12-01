package entity

import (
	"time"
)

// category.go
type Category struct {
	CategoryID   string `gorm:"primaryKey;autoIncrement:false"`
	UserID       string `gorm:"index"` // แค่ index พอ ไม่ต้อง composite PK
	CategoryName string `gorm:"length:20"`
	Budget       float64
	ColorCode    string    `gorm:"length:7"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	// FK ไปหา User ทางเดียว
	User User `gorm:"foreignKey:UserID;references:UserID"`

	// เอา constraint ออกไว้ก่อน ให้มีไว้แค่ preload relationship
	Transactions []Transaction `gorm:"foreignKey:UserID,CategoryID;references:UserID,CategoryID"`
}
