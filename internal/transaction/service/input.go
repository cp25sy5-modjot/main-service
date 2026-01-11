package transactionsvc

import (
	"time"
)

type TransactionCreateInput struct {
	Date  time.Time
	Items []TransactionItemInput
}

type TransactionItemInput struct {
	Title      string
	Price      float64
	CategoryID string
}

type TransactionUpdateInput struct {
	Date  *time.Time
	Items []TransactionItemInput
}
