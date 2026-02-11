package draft

import "time"

type DraftStatus string

const (
	DraftStatusProcessing     DraftStatus = "processing"
	DraftStatusWaitingConfirm DraftStatus = "waiting_confirm"
	DraftStatusFailed         DraftStatus = "failed"
)

type DraftTxn struct {
	TraceID string `json:"trace_id"`
	UserID  string `json:"user_id"`

	Status DraftStatus `json:"status"` // processing | waiting_confirm | failed
	Title  string      `json:"title,omitempty"`
	Date   time.Time   `json:"date,omitempty"`
	Items  []DraftItem `json:"items,omitempty"`

	Error string `json:"error,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type DraftItem struct {
	Title      string  `json:"title"`
	Price      float64 `json:"price"`
	CategoryID string  `json:"category_id"`
}

type NewDraftRequest struct {
	Title     string      `json:"title"`
	Date      time.Time   `json:"date"`
	Items     []DraftItem `json:"items" validate:"required,min=1,dive"`
	CreatedAt time.Time   `json:"created_at,omitempty"`
}

type ConfirmRequest struct {
	Title string      `json:"title"`
	Date  *time.Time  `json:"date"`
	Items []DraftItem `json:"items" validate:"required,min=1,dive"`
}

type UpdateDraftStatusRequest struct {
	Status DraftStatus `json:"status"`
	Error  string      `json:"error,omitempty"`
}

type DraftStats struct {
	Total          int `json:"total"`
	Processing     int `json:"processing"`
	WaitingConfirm int `json:"waiting_confirm"`
	Failed         int `json:"failed"`
}

type DraftRes struct {
	TraceID string `json:"trace_id"`

	Status DraftStatus    `json:"status"` // processing | waiting_confirm | failed
	Title  string         `json:"title,omitempty"`
	Date   time.Time      `json:"date,omitempty"`
	Items  []DraftItemRes `json:"items,omitempty"`

	Error string `json:"error,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

type DraftItemRes struct {
	Title         string  `json:"title"`
	Price         float64 `json:"price"`
	CategoryID    string  `json:"category_id"`
	CategoryName  string  `json:"category_name,omitempty"`
	CategoryIcon  string  `json:"category_icon,omitempty"`
	CategoryColor string  `json:"category_color,omitempty"`
}
