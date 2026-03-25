package model

import (
	"time"
)

type FixCostCreateReq struct {
	Title      string  `json:"title"`
	Price      float64 `json:"price"`
	CategoryID string  `json:"category_id"`

	StartDate     time.Time  `json:"start_date"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	RemainingRuns *int       `json:"remaining_runs,omitempty"`

	IntervalType  string `json:"interval_type" validate:"required,oneof=daily weekly monthly yearly"` // daily, weekly, monthly, yearly
	IntervalValue int    `json:"interval_value" validate:"required,min=1"`                            // e.g., every 2 weeks
}

type FixCostCreateInput struct {
	UserID     string
	Title      string
	Price      float64
	CategoryID string

	StartDate     time.Time
	EndDate       *time.Time
	RemainingRuns *int

	IntervalType  string
	IntervalValue int
}

type FixCostUpdateReq struct {
	Title      *string  `json:"title"`
	Price      *float64 `json:"price"`
	CategoryID *string  `json:"category_id"`

	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	RemainingRuns *int       `json:"remaining_runs"`
	Status        *string    `json:"status" validate:"oneof=active paused finished"`

	IntervalType  *string `json:"interval_type" validate:"oneof=daily weekly monthly yearly"` // daily, weekly, monthly, yearly
	IntervalValue *int    `json:"interval_value" validate:"min=1"`                            // e.g., every 2 weeks
}

type FixCostUpdateInput struct {
	FixCostID string
	UserID    string

	Title      *string
	Price      *float64
	CategoryID *string

	StartDate     *time.Time
	EndDate       *time.Time
	RemainingRuns *int
	Status        *string

	IntervalType  *string
	IntervalValue *int
}

type FixCostRes struct {
	FixCostID  string  `json:"fix_cost_id"`
	Title      string  `json:"title"`
	Price      float64 `json:"price"`
	CategoryID string  `json:"category_id"`

	StartDate     time.Time  `json:"start_date"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	RemainingRuns *int       `json:"remaining_runs,omitempty"`

	IntervalType  string     `json:"interval_type"`  // daily, weekly, monthly, yearly
	IntervalValue int        `json:"interval_value"` // e.g., every 2 weeks
	NextRunDate   time.Time  `json:"next_run_date"`
	LastRunAt     *time.Time `json:"last_run_at,omitempty"`
	
	Status        string     `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	CategoryName  string `json:"category_name"`
	CategoryIcon  string `json:"category_icon"`
	CategoryColor string `json:"category_color"`
}
