package transactionsvc

import (
	"time"
)

type TransactionCreateInput struct {
	Title      string
	Price      float64
	Quantity   float64
	CategoryID *string
	Date       time.Time
}

type TransactionUpdateInput struct {
	Title      string
	Price      float64
	Quantity   float64
	CategoryID *string
	Date       time.Time
}
