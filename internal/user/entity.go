package user

import "time"

type User struct {
	UserID    string    `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	Name      string    `gorm:"length:100" json:"name"`
	DOB       time.Time `json:"dob"`
	Email     string    `gorm:"length:100,unique" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
