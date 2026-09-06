package models

import (
	"habitask-backend-go/lib/constants"
	"time"
)

type Task struct {
	BaseModel
	Name              string             `gorm:"unique;varchar(100)"`
	Description       string             `gorm:"Varchar(500)"`
	UserID            uint               `gorm:"index"`
	User              *User              `gorm:"foreignKey:UserID"`
	Category          constants.Category `gorm:"varchar(50)"`
	IsCompleted       *bool              `gorm:"default:false"`
	CompletedAt       *time.Time         `gorm:"datetime"`
	DueAt             *time.Time         `gorm:"datetime"`
	EstimatedDuration int                `gorm:"default:0"`
	ActualDuration    int                `gorm:"default:0"`
	DurationUnit      string             `gorm:"enum('seconds', 'minutes', 'hours', 'days', 'weeks', 'months', 'years');Varchar(50)"`
}
