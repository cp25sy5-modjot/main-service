package entity

import "time"

type FixCostStatus string

const (
	FixCostStatusActive   FixCostStatus = "active"
	FixCostStatusPaused   FixCostStatus = "paused"
	FixCostStatusFinished FixCostStatus = "finished"
)

type IntervalType string

const (
	IntervalDaily   IntervalType = "daily"
	IntervalWeekly  IntervalType = "weekly"
	IntervalMonthly IntervalType = "monthly"
	IntervalYearly  IntervalType = "yearly"
)

type FixCost struct {
	FixCostID  string `gorm:"primaryKey;autoIncrement:false"`
	UserID     string
	Title      string
	Price      float64
	CategoryID string

	StartDate time.Time
	EndDate   *time.Time
	MaxRun    *int // Number of maximum runs (optional, nil means unlimited)
	RunCount  int  // Number of times the fix cost has been processed

	IntervalType  IntervalType // "daily", "weekly", "monthly", "yearly"
	IntervalValue int          // Number of intervals between runs (e.g., every 2 weeks)

	LastRunAt   *time.Time    // Last time the fix cost was processed
	NextRunDate time.Time     `gorm:"index"`
	Status      FixCostStatus `gorm:"index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID"`
}

func (fc *FixCost) IsActive() bool {
	now := time.Now()

	if fc.EndDate != nil && now.After(*fc.EndDate) {
		return false
	}
	if fc.MaxRun != nil && *fc.MaxRun <= 0 {
		return false
	}
	return true
}
