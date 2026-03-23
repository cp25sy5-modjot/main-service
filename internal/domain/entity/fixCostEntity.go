package entity

import "time"

type FixCostStatus string

const (
	FixCostStatusActive    FixCostStatus = "active"
	FixCostStatusPaused    FixCostStatus = "paused"
	FixCostStatusCompleted FixCostStatus = "finished"
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

	StartDate     time.Time
	EndDate       *time.Time
	RemainingRuns *int // Number of remaining runs, if applicable

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
