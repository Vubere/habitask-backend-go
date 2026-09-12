package models

import (
	"habitask-backend-go/pkg/lib/constants"
	"time"
)

type Task struct {
	BaseModel
	Name              string             `gorm:"unique;varchar(100)"`
	Description       string             `gorm:"Varchar(500)"`
	UserID            string             `gorm:"index;type:varchar(36)"`
	User              *User              `gorm:"foreignKey:UserID"`
	Category          constants.Category `gorm:"varchar(50)"`
	IsCompleted       *bool              `gorm:"default:false"`
	CompletedAt       *time.Time         `gorm:"datetime"`
	DueAt             *time.Time         `gorm:"datetime"`
	EstimatedDuration int                `gorm:"default:0"` //seconds
	ActualDuration    int                `gorm:"default:0"` //seconds
	Priority          string             `gorm:"enum('low', 'medium', 'high');Varchar(50)"`
}

type TaskSummary struct {
	Label           string  `json:"label"`
	Count           int     `json:"count"`
	CompletedCount  int     `json:"completed_count"`
	IncompleteCount int     `json:"incomplete_count"`
	DueCount        int     `json:"due_count"`
	LastTaskDate    *string `json:"last_done"`
}

type TaskQuery struct {
	Task
	Search               string     `json:"search"`
	CompletedAtLte       *time.Time `json:"completed_at_lte"`
	CompletedAtGte       *time.Time `json:"completed_at_gte"`
	DueAtLte             *time.Time `json:"due_at_lte"`
	DueAtGte             *time.Time `json:"due_at_gte"`
	EstimatedDurationLte *int       `json:"estimated_duration_lte"`
	EstimatedDurationGte *int       `json:"estimated_duration_gte"`
	ActualDurationLte    *int       `json:"actual_duration_lte"`
	ActualDurationGte    *int       `json:"actual_duration_gte"`
	//SummaryQuerys
	GroupBy   string `json:"group_by"`
	DateGroup string `json:"date_group"` //day, week, month, year, week_day, month_name
}
