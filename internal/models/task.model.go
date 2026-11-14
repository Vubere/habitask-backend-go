package models

import "time"

type TaskStatus string

const (
	TaskStatusTodo TaskStatus = "todo"
	TaskStatusDone TaskStatus = "done"
)

type Task struct {
	BaseModel
	Title       string     `gorm:"varchar(100);not null" json:"title"`
	Notes       string     `gorm:"varchar(500)" json:"notes"`
	Icon        string     `gorm:"varchar(50)" json:"icon"`
	Color       string     `gorm:"varchar(20)" json:"color"`
	ArchivedAt  *time.Time `gorm:"datetime" json:"archived_at,omitempty"`
	UserID      string     `gorm:"index;varchar(36)" json:"user_id"`
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status      TaskStatus `gorm:"varchar(10);default:todo" json:"status"`
	DueDate     *time.Time `gorm:"type:date" json:"due_date,omitempty"`
	DueTime     *string    `gorm:"varchar(5)" json:"due_time,omitempty"` // "HH:mm"
	Priority    string     `gorm:"varchar(10);default:none" json:"priority"`
	Subtasks    string     `gorm:"type:json" json:"subtasks"` // JSON []Subtask
	CompletedAt *time.Time `gorm:"datetime" json:"completed_at,omitempty"`
}

type TaskSummary struct {
	Label           string  `json:"label"`
	Count           int     `json:"count"`
	CompletedCount  int     `json:"completed_count"`
	IncompleteCount int     `json:"incomplete_count"`
	DueCount        int     `json:"due_count"`
	LastTaskDate    *string `json:"last_task_date"`
}

type TaskQuery struct {
	Task
	Search   string     `json:"search"`
	DateLte  *time.Time `json:"date_lte"`
	DateGte  *time.Time `json:"date_gte"`
	Archived *bool      `json:"archived"`
	//SummaryQuerys
	GroupBy   string `json:"group_by"`
	DateGroup string `json:"date_group"` //day, week, month, year, week_day, month_name
}
