package model

import (
	"time"
)

type TransactionSearchParams struct {
	TransactionID string `json:"transaction_id"`
	UserID        string `json:"user_id" validate:"required"`
}
type TransactionItemReq struct {
	Title      string  `json:"title" validate:"required,min=2,max=20"`
	Price      float64 `json:"price" validate:"required"`
	CategoryID string  `json:"category_id" validate:"required"`
}

type TransactionInsertReq struct {
	Title string               `json:"title" validate:"required,min=2,max=20"`
	Date  time.Time            `json:"date"`
	Items []TransactionItemReq `json:"items" validate:"required,min=1,dive"`
}

type TransactionUpdateReq struct {
	Title string               `json:"title" validate:"required,min=2,max=20"`
	Date  *time.Time           `json:"date" validate:"required"`
	Items []TransactionItemReq `json:"items" validate:"required,min=1,dive"`
}

type TransactionRes struct {
	TransactionID string               `json:"transaction_id"`
	Date          time.Time            `json:"date"`
	Type          string               `json:"type"`
	Total         float64              `json:"total"`
	Items         []TransactionItemRes `json:"items"`
}

type TransactionFilter struct {
	Date *time.Time `json:"date"`
}

type TransactionCompareMonthResponse struct {
	Transactions       []TransactionRes `json:"transactions"`
	CurrentMonthTotal  float64          `json:"current_month_total"`
	PreviousMonthTotal float64          `json:"previous_month_total"`
}
