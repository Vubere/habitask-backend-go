package models

import "time"

type TaskSteps struct {
	BaseModel
	TaskID            string     `gorm:"index"`
	Task              *Task      `gorm:"foreignKey:TaskID"`
	User              *User      `gorm:"foreignKey:UserID"`
	UserID            string     `gorm:"index"`
	Description       string     `gorm:"Varchar(500)"`
	IsDone            *bool      `gorm:"default:false"`
	DoneAt            *time.Time `gorm:"datetime"`
	TimeDue           *time.Time `gorm:"datetime"`
	EstimatedDuration int        `gorm:"default:0"`
	ActualDuration    int        `gorm:"default:0"`
	DurationUnit      string     `gorm:"enum('seconds', 'minutes', 'hours', 'days', 'weeks', 'months', 'years');Varchar(50)"`
	Priority          string     `gorm:"enum('low', 'medium', 'high');Varchar(50)"`
}
