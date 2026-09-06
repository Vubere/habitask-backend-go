package models

import (
	"habitask-backend-go/pkg/lib/constants"
	"time"
)

type Habit struct {
	BaseModel
	Name        string             `gorm:"unique;varchar(100)"`
	Description string             `gorm:"varchar(500)"`
	UserID      string             `gorm:"index"`
	User        *User              `gorm:"foreignKey:UserID"`
	Category    constants.Category `gorm:"varchar(255)"`
	Pros        string             `gorm:"text"`
	Cons        string             `gorm:"text"`
	IsPositive  *bool              `gorm:"default:true"`
	LastDone    *time.Time         `gorm:"datetime"`
}

type HabbitSummary struct {
	Label          string  `json:"label"`
	Count          int     `json:"count"`
	PositiveCount  int     `json:"positive_count"`
	NegativeCount  int     `json:"negative_count"`
	LastHabbitDate *string `json:"last_done"`
}

type HabitQuery struct {
	Habit
	Search      string     `json:"search"`
	LastDoneLte *time.Time `json:"last_done_lte"`
	LastDoneGte *time.Time `json:"last_done_gte"`
	//SummaryFiltes
	GroupBy   string `json:"group_by"`
	DateGroup string `json:"date_group"` //DAY, WEEK, MONTH, YEAR
}
