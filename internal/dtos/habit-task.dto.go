package dtos

import "time"

type HabitTaskDTO struct {
	Id          int    `json:"id" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Category    string `json:"category" gorm:"column:category"`
	Status      string `json:"status" gorm:"column:status"`
	Progress    int    `json:"progress" gorm:"column:progress"`
	Streak      *int   `json:"streak" gorm:"column:streak"`
	CreatedAt   string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   string `json:"updated_at" gorm:"column:updated_at"`
}

type HabitTaskFilterDTO struct {
	CreatedAtLTE *time.Time `json:"created_at_lte"`
	CreatedAtGTE *time.Time `json:"created_at_gte"`
	UserID       string     `json:"user_id"`
}
