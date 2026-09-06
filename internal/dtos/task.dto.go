package dtos

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/lib/constants"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskDTO struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Category          constants.Category `json:"category"`
	Description       string             `json:"description"`
	IsCompleted       *bool              `json:"is_completed"`
	CompletedAt       string             `json:"completed_at"`
	DueAt             string             `json:"due_at"`
	EstimatedDuration int                `json:"estimated_duration"`
	ActualDuration    int                `json:"actual_duration"`
	DurationUnit      string             `json:"duration_unit"`
}

func ParseTask(task *models.Task) *TaskDTO {
	return &TaskDTO{
		ID:                task.ID,
		Name:              task.Name,
		Category:          task.Category,
		Description:       task.Description,
		IsCompleted:       task.IsCompleted,
		CompletedAt:       task.CompletedAt.Format("2006-01-02 15:04:05"),
		DueAt:             task.DueAt.Format("2006-01-02 15:04:05"),
		EstimatedDuration: task.EstimatedDuration,
		ActualDuration:    task.ActualDuration,
		DurationUnit:      task.DurationUnit,
	}
}

type TaskCreateDTO struct {
	Name              string             `json:"name" binding:"required"`
	Category          constants.Category `json:"category" binding:"required"`
	Description       string             `json:"description" binding:"required"`
	EstimatedDuration int                `json:"estimated_duration"`
	DurationUnit      string             `json:"duration_unit"`
}

func (t *TaskCreateDTO) ToTask() *models.Task {
	task := &models.Task{
		Name:              t.Name,
		Category:          t.Category,
		Description:       t.Description,
		EstimatedDuration: t.EstimatedDuration,
		DurationUnit:      t.DurationUnit,
	}
	return task
}

type TaskUpdateDTO struct {
	Name              *string             `json:"name"`
	Category          *constants.Category `json:"category"`
	Description       *string             `json:"description"`
	EstimatedDuration *int                `json:"estimated_duration"`
	DurationUnit      *string             `json:"duration_unit"`
	ActualDuration    *int                `json:"actual_duration"`
}

func (t *TaskUpdateDTO) ToTask() *models.Task {
	task := &models.Task{}
	if t.Name != nil {
		task.Name = *t.Name
	}
	if t.Category != nil {
		task.Category = *t.Category
	}
	if t.Description != nil {
		task.Description = *t.Description
	}
	return task
}

type TaskCompleteDTO struct {
	CompletedAt string `json:"completed_at"`
}

func (t *TaskCompleteDTO) ToTask(ctx *gin.Context) *models.Task {
	task := &models.Task{}
	if t.CompletedAt != "" {
		ca, err := time.Parse("2006-01-02 15:04:05", t.CompletedAt)
		if err != nil {
			if ctx != nil {
				ctx.JSON(400, gin.H{"error": "invalid date format on field completed_at"})
			}
			return nil
		}
		task.CompletedAt = &ca
	}
	boolTrue := true
	task.IsCompleted = &boolTrue
	return task
}
