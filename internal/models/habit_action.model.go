package models

import "habitask-backend-go/lib/constants"

type HabitAction struct {
	BaseModel
	Title        string             `gorm:"unique;varchar(100)"`
	Description  string             `gorm:"Varchar(500)"`
	HabitID      string             `gorm:"index"`
	Habit        *Habit             `gorm:"foreignKey:HabitID"`
	CostIncurred float64            `gorm:"default:0"`
	Category     constants.Category `gorm:"varchar(255)"`
	IsPositive   *bool              `gorm:"default:true"`
	UserID       string             `gorm:"index"`
	User         *User              `gorm:"foreignKey:UserID"`
}
