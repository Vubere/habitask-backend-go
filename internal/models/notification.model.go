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
	UserID      string    `gorm:"index;varchar(36)"`
	User        *User     `gorm:"foreignKey:UserID"`
}

type NotificationSummary struct {
	Label                string  `json:"label"`
	Count                int     `json:"count"`
	LastNotificationDate *string `json:"last_notification_date"`
	UnreadCount          int     `json:"unread_count"`
}

type NotificationQuery struct {
	Notification
	Search  string  `json:"search"`
	DateLte *string `json:"date_lte"`
	DateGte *string `json:"date_gte"`
	//SummaryQuerys
	GroupBy   string `json:"group_by"`
	DateGroup string `json:"date_group"` //day, week, month, year, week_day, month_name
}
