package model

type ExpenseSummaryRes struct {
	Period string           `json:"period"`
	Data   []ExpenseSummary `json:"data"`
}

type ExpenseSummary struct {
	Key   string  `json:"key"`
	Label string  `json:"label"`
	Total float64 `json:"total"`
}

type CategorySummaryRes struct {
	Period  string            `json:"period"`
	Total   float64           `json:"total"`
	Average float64           `json:"average"`
	Data    []CategorySummary `json:"data"`
}

type CategorySummary struct {
	CategoryID    string  `json:"category_id"`
	CategoryIcon  string  `json:"category_icon"`
	CategoryName  string  `json:"category_name"`
	CategoryColor string  `json:"category_color"`
	Total         float64 `json:"total"`
}

type CategorySummaryQuery struct {
	Period          string `query:"period" validate:"required,oneof=day week month year past_year"`
	Date            string `query:"date" validate:"required"`
	IsUseDateFilter bool   `query:"is_use_date"`
}

type ExpenseSummaryQuery struct {
	Period          string `query:"period" validate:"required,oneof=week month year"`
	Date            string `query:"date" validate:"required"`
	IsUseDateFilter bool   `query:"is_use_date"`
}
