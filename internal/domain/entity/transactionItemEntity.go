package entity

type TransactionItem struct {
	TransactionID string `gorm:"primaryKey"`
	ItemID        string `gorm:"primaryKey"`

	Title      string
	Price      float64
	CategoryID string

	// Relationships
	Category    Category    `gorm:"foreignKey:CategoryID;references:CategoryID"`
	Transaction Transaction `gorm:"foreignKey:TransactionID;references:TransactionID"`
}
