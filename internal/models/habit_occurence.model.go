package models

import "time"

type HabitOccurence struct {
	BaseModel
	HabitID string    `gorm:"index"`
	Habit   *Habit    `gorm:"foreignKey:HabitID"`
	UserID  string    `gorm:"index"`
	User    *User     `gorm:"foreignKey:UserID"`
	Date    time.Time `gorm:"datetime"`
}
