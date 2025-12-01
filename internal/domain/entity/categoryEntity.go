package entity

import (
	"time"
)

type Category struct {
	CategoryID   string    `gorm:"primaryKey;autoIncrement:false"`
	UserID       string    `gorm:"primaryKey;autoIncrement:false"`
	CategoryName string    `gorm:"length:20"`
	Budget       float64
	ColorCode    string    `gorm:"length:7"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	// belongs-to User
	User User `gorm:"foreignKey:UserID;references:UserID"`

	// has-many Transactions – ลบ category แล้วลบทุก transaction ใน category นี้
	Transactions []Transaction `gorm:"foreignKey:UserID,CategoryID;references:UserID,CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}


