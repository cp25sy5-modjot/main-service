package model

import "time"

type LastTransaction struct {
	TransactionID     string    `json:"transaction_id"`
	ItemID            string    `json:"item_id"`
	Title             string    `json:"title"`
	Price             float64   `json:"price"`
	Date              time.Time `json:"date"`
	Type              string    `json:"type"`
	CategoryID        *string   `json:"category_id"`
	CategoryName      string    `json:"category_name"`
	CategoryColorCode string    `json:"category_color_code"`
}

type TopCategoryUsage struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	ColorCode    string  `json:"color_code"`
	Budget       float64 `json:"budget"`
	BudgetUsage  float64 `json:"budget_usage"`
}

type OverviewResponse struct {
	LastTransactions []LastTransaction  `json:"last_transactions"`
	TopCategories    []TopCategoryUsage `json:"top_categories"`
}
