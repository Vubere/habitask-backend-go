package models

import "time"

type TaskSteps struct {
	BaseModel
	TaskID            string     `gorm:"index;varchar(36)"`
	Task              *Task      `gorm:"foreignKey:TaskID"`
	User              *User      `gorm:"foreignKey:UserID"`
	UserID            string     `gorm:"index;type:varchar(36)"`
	Description       string     `gorm:"Varchar(500)"`
	IsDone            *bool      `gorm:"default:false"`
	DoneAt            *time.Time `gorm:"datetime"`
	TimeDue           *time.Time `gorm:"datetime"`
	EstimatedDuration int        `gorm:"default:0"` //seconds
	ActualDuration    int        `gorm:"default:0"` //seconds
	Order             *int       `gorm:"default:0"`
}

type TaskStepsQuery struct {
	TaskSteps
	Search               string     `json:"search"`
	DoneAtLte            *time.Time `json:"done_at_lte"`
	DoneAtGte            *time.Time `json:"done_at_gte"`
	TimeDueLte           *time.Time `json:"time_due_lte"`
	TimeDueGte           *time.Time `json:"time_due_gte"`
	EstimatedDurationLte *int       `json:"estimated_duration_lte"`
	EstimatedDurationGte *int       `json:"estimated_duration_gte"`
	ActualDurationLte    *int       `json:"actual_duration_lte"`
	ActualDurationGte    *int       `json:"actual_duration_gte"`
	//SummaryQuerys
	GroupBy   string `json:"group_by"`
	DateGroup string `json:"date_group"` //day, week, month, year, week_day, month_name
}

type TaskStepSummary struct {
	Label         string  `json:"label"`
	Count         int     `json:"count"`
	DoneCount     int     `json:"done_count"`
	LastTaskStep  *string `json:"last_task_step"`
	TimeDue       *string `json:"time_due"`
	EstimatedTime int     `json:"estimated_time"`
	ActualTime    int     `json:"actual_time"`
	TimeUnit      string  `json:"time_unit"`
}
