package models

import (
	"habitask-backend-go/lib/constants"
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
