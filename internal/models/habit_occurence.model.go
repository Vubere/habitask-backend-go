package models

import "time"

type HabitOccurence struct {
	BaseModel
	HabitID string    `gorm:"index"`
	Habit   *Habit    `gorm:"foreignKey:HabitID"`
	UserID  string    `gorm:"index;varchar(36)"`
	User    *User     `gorm:"foreignKey:UserID"`
	Date    time.Time `gorm:"datetime"`
}

type HabitOccurenceQuery struct {
	HabitOccurence
	Search  string  `json:"search"`
	DateLte *string `json:"date_lte"`
	DateGte *string `json:"date_gte"`
	//SummaryQuerys
	GroupBy   string `json:"group_by"`
	DateGroup string `json:"date_group"` //day, week, month, year, week_day, month_name
}

type HabitOccurenceSummary struct {
	Label         string  `json:"label"`
	Count         int     `json:"count"`
	LastHabitDate *string `json:"last_habit_date"`
	PositiveCount int     `json:"positive_count"`
	NegativeCount int     `json:"negative_count"`
	UserID        string  `json:"user_id"`
}
