package models

import "time"

type HabitCompletion struct {
	BaseModel
	HabitID     string    `gorm:"index;varchar(36);not null" json:"habit_id"`
	Habit       *Habit    `gorm:"foreignKey:HabitID" json:"habit,omitempty"`
	UserID      string    `gorm:"index;varchar(36)" json:"user_id"`
	Date        time.Time `gorm:"type:date;not null" json:"date"`
	Count       int       `gorm:"default:0" json:"count"`
	CompletedAt time.Time `gorm:"datetime" json:"completed_at"`
}

type HabitCompletionSummary struct {
	Label         string  `json:"label"`
	Count         int     `json:"count"`
	TotalCount    int     `json:"total_count"`
	LastHabitDate *string `json:"last_habit_date"`
}

type HabitCompletionQuery struct {
	HabitCompletion
	Search  string     `json:"search"`
	DateLte *time.Time `json:"date_lte"`
	DateGte *time.Time `json:"date_gte"`
	//SummaryQuerys
	GroupBy   string `json:"group_by"`
	DateGroup string `json:"date_group"` //day, week, month, year, week_day, month_name
}
