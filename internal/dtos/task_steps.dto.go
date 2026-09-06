package dtos

import (
	"habitask-backend-go/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskStepDTO struct {
	ID            string     `json:"id"`
	Description   string     `json:"description"`
	IsDone        *bool      `json:"is_done"`
	DoneAt        *time.Time `json:"done_at"`
	TimeDue       *time.Time `json:"time_due"`
	EstimatedTime int        `json:"estimated_time"`
	ActualTime    int        `json:"actual_time"`
	TimeUnit      string     `json:"time_unit"`
	Priority      string     `json:"priority"`
	User          *UserDTO   `json:"user"`
	Task          *TaskDTO   `json:"task"`
}

func ParseTaskStep(taskStep *models.TaskSteps) *TaskStepDTO {
	return &TaskStepDTO{
		ID:            taskStep.ID,
		Description:   taskStep.Description,
		IsDone:        taskStep.IsDone,
		DoneAt:        taskStep.DoneAt,
		TimeDue:       taskStep.TimeDue,
		EstimatedTime: taskStep.EstimatedDuration,
		ActualTime:    taskStep.ActualDuration,
		TimeUnit:      taskStep.DurationUnit,
		Priority:      taskStep.Priority,
		User:          ParseUser(taskStep.User),
		Task:          ParseTask(taskStep.Task),
	}
}

type TaskStepCreateDTO struct {
	TaskID            string `json:"task_id" binding:"required"`
	Description       string `json:"description" binding:"required"`
	EstimatedDuration int    `json:"estimated_time" binding:"required"`
	TimeDue           string `json:"time_due"`
	DurationUnit      string `json:"time_unit"`
	Priority          string `json:"priority"`
}

func (t *TaskStepCreateDTO) ToTaskStep(ctx *gin.Context) *models.TaskSteps {
	taskStep := &models.TaskSteps{
		Description:       t.Description,
		EstimatedDuration: t.EstimatedDuration,
		DurationUnit:      t.DurationUnit,
		Priority:          t.Priority,
		TaskID:            t.TaskID,
	}
	if t.TimeDue != "" {
		td, err := time.Parse("2006-01-02 15:04:05", t.TimeDue)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid date format"})
			return nil
		}
		taskStep.TimeDue = &td
	}
	if t.DurationUnit != "" {
		taskStep.DurationUnit = t.DurationUnit
	}
	return taskStep
}

type TaskStepUpdateDTO struct {
	Description       *string `json:"description"`
	EstimatedDuration *int    `json:"estimated_duration"`
	TimeDue           *string `json:"time_due"`
	DurationUnit      *string `json:"time_unit"`
	Priority          *string `json:"priority"`
	ActualDuration    *int    `json:"actual_duration"`
}

func (t *TaskStepUpdateDTO) ToTaskStep(ctx *gin.Context) *models.TaskSteps {
	taskStep := &models.TaskSteps{}
	if t.Description != nil {
		taskStep.Description = *t.Description
	}
	if t.EstimatedDuration != nil {
		taskStep.EstimatedDuration = *t.EstimatedDuration
	}
	if t.TimeDue != nil {
		timeDue, err := time.Parse("2006-01-02 15:04:05", *t.TimeDue)
		if err != nil {
			if ctx != nil {
				ctx.JSON(400, gin.H{"error": "invalid date format on field time_due"})
			}
			return nil
		}
		taskStep.TimeDue = &timeDue
	}
	if t.DurationUnit != nil {
		taskStep.DurationUnit = *t.DurationUnit
	}
	if t.Priority != nil {
		taskStep.Priority = *t.Priority
	}
	return taskStep
}

type TaskStepCompleteDTO struct {
	DoneAt string `json:"done_at"`
}

func (t *TaskStepCompleteDTO) ToTaskStep(ctx *gin.Context) *models.TaskSteps {
	taskStep := &models.TaskSteps{}
	if t.DoneAt != "" {
		da, err := time.Parse("2006-01-02 15:04:05", t.DoneAt)
		if err != nil {
			if ctx != nil {
				ctx.JSON(400, gin.H{"error": "invalid date format on field done_at"})
			}
			return nil
		}
		taskStep.DoneAt = &da
	}
	boolTrue := true
	taskStep.IsDone = &boolTrue
	return taskStep
}
