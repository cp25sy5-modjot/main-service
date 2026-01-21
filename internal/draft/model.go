package draft

import "time"

type DraftStatus string

const (
	DraftStatusProcessing      DraftStatus = "processing"
	DraftStatusWaitingConfirm  DraftStatus = "waiting_confirm"
	DraftStatusFailed          DraftStatus = "failed"
)

type DraftTxn struct {
	TraceID string `json:"trace_id"`
	UserID  string `json:"user_id"`

	Status DraftStatus `json:"status"` // processing | waiting_confirm | failed

	Date  time.Time   `json:"date,omitempty"`
	Items []DraftItem `json:"items,omitempty"`

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
	Date  time.Time   `json:"date"`
	Items []DraftItem  `json:"items" validate:"required,min=1,dive"`
}

type ConfirmRequest struct {
	Date  *time.Time   `json:"date"`
	Items []DraftItem  `json:"items" validate:"required,min=1,dive"`
}
