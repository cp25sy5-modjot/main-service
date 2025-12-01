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

	// belongs-to User (ให้มี FK จาก categories.user_id -> users.user_id แค่ทางเดียว)
	User User `gorm:"foreignKey:UserID;references:UserID"`

	// ถ้าไม่อยากให้ GORM สร้าง FK/constraint เพิ่มบน users/transactions
	// เอา constraint ออก เหลือแค่ preload relationship เฉย ๆ ก็ได้
	Transactions []Transaction `gorm:"foreignKey:UserID,CategoryID;references:UserID,CategoryID"`
}



