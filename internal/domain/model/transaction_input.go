package model

import (
	"time"
)

type TransactionCreateInput struct {
	title string
	Date  time.Time
	Items []TransactionItemInput
}

type TransactionItemInput struct {
	Title      string
	Price      float64
	CategoryID string
}

type TransactionUpdateInput struct {
	Title string
	Date  *time.Time
	Items []TransactionItemInput
}
