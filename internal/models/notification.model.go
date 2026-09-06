package models

import "time"

type Notification struct {
	BaseModel
	Title       string    `gorm:"unique;varchar(100)"`
	Description string    `gorm:"varchar(500)"`
	Category    string    `gorm:"enum('reminder', 'notification');Varchar(50)"`
	Date        time.Time `gorm:"datetime"`
	IsRead      *bool     `gorm:"default:false"`
	Source      string    `gorm:"enum('habit', 'habit-action', 'task', 'task_step');Varchar(50)"`
	SourceID    string    `gorm:"index"`
	UserID      string    `gorm:"index"`
	User        *User     `gorm:"foreignKey:UserID"`
}
