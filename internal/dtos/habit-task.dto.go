package dtos

import "time"

type HabitTaskDTO struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Notes     string `json:"notes"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	Target    int    `json:"target"`
	Streak    *int   `json:"streak"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type HabitTaskFilterDTO struct {
	CreatedAtLTE *time.Time `json:"created_at_lte"`
	CreatedAtGTE *time.Time `json:"created_at_gte"`
	Type         string     `json:"type"`
	UserID       string     `json:"user_id"`
}
