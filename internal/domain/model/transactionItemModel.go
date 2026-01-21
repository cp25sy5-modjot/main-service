package model

type TransactionItemInsertReq struct {
	Title      string  `json:"title" validate:"required,min=2,max=20"`
	Price      float64 `json:"price" validate:"required"`
	CategoryID string  `json:"category_id" validate:"required"`
}

type TransactionItemUpdateReq struct {
	Title      string  `json:"title" validate:"required,min=2,max=20"`
	Price      float64 `json:"price" validate:"required"`
	CategoryID string  `json:"category_id" validate:"required"`
}

type TransactionItemRes struct {
	TransactionID     string  `json:"transaction_id"`
	ItemID            string  `json:"item_id"`
	Title             string  `json:"title"`
	Price             float64 `json:"price"`
	CategoryID        string  `json:"category_id"`
	CategoryName      string  `json:"category_name"`
	CategoryColorCode string  `json:"category_color_code"`
	Icon              string  `json:"icon"`
}

type TransactionItemSearchParams struct {
	TransactionID string `json:"transaction_id"`
	UserID        string `json:"user_id" validate:"required"`
	ItemID        string `json:"item_id"`
}
