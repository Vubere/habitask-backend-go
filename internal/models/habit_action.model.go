package models

import "habitask-backend-go/pkg/lib/constants"

type HabitAction struct {
	BaseModel
	Title        string             `gorm:"unique;varchar(100)"`
	Description  string             `gorm:"varchar(500)"`
	HabitID      string             `gorm:"index;varchar(36)"`
	Habit        *Habit             `gorm:"foreignKey:HabitID"`
	CostIncurred float64            `gorm:"default:0"`
	Category     constants.Category `gorm:"varchar(255)"`
	IsPositive   *bool              `gorm:"default:true"`
	UserID       string             `gorm:"index;type:varchar(36)"`
	User         *User              `gorm:"foreignKey:UserID"`
}

type HabitActionSummary struct {
	Label        string  `json:"label"`
	Count        int     `json:"count"`
	CostIncurred float64 `json:"cost_incurred"`
}

type HabitActionQuery struct {
	HabitAction
	Search  string  `json:"search"`
	DateLte *string `json:"date_lte"`
	DateGte *string `json:"date_gte"`
	//SummaryQuerys
	GroupBy   string `json:"group_by"`
	DateGroup string `json:"date_group"` //day, week, month, year, week_day, month_name
}
