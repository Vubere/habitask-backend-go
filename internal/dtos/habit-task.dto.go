package dtos

import "time"

type HabitTaskDTO struct {
	Id          string `json:"id" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Category    string `json:"category" gorm:"column:category"`
	Status      string `json:"status" gorm:"column:status"`
	Progress    int    `json:"progress" gorm:"column:progress"`
	Steps       int    `json:"steps" gorm:"column:steps"`
	Streak      *int   `json:"streak" gorm:"column:streak"`
	Type        string `json:"type" gorm:"column:type"`
	CreatedAt   string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   string `json:"updated_at" gorm:"column:updated_at"`
}

type HabitTaskFilterDTO struct {
	CreatedAtLTE *time.Time `json:"created_at_lte"`
	CreatedAtGTE *time.Time `json:"created_at_gte"`
	Type         string     `json:"type"`
	UserID       string     `json:"user_id"`
}
