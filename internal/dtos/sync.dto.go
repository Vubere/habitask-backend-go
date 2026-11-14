package dtos

import (
	"habitask-backend-go/internal/models"
	"time"
)

type SyncHabitDTO struct {
	HabitDTO
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type SyncTaskDTO struct {
	TaskDTO
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type SyncHabitCompletionDTO struct {
	HabitCompletionDTO
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type SyncDTO struct {
	Habits           []SyncHabitDTO           `json:"habits"`
	Tasks            []SyncTaskDTO            `json:"tasks"`
	HabitCompletions []SyncHabitCompletionDTO `json:"habit_completions"`
	ServerTime       time.Time                `json:"server_time"`
}

func ParseSyncHabit(habit *models.Habit) SyncHabitDTO {
	dto := SyncHabitDTO{}
	if parsed := ParseHabit(habit); parsed != nil {
		dto.HabitDTO = *parsed
	}
	if habit != nil && habit.DeletedAt.Valid {
		deletedAt := habit.DeletedAt.Time
		dto.DeletedAt = &deletedAt
	}
	return dto
}

func ParseSyncTask(task *models.Task) SyncTaskDTO {
	dto := SyncTaskDTO{}
	if parsed := ParseTask(task); parsed != nil {
		dto.TaskDTO = *parsed
	}
	if task != nil && task.DeletedAt.Valid {
		deletedAt := task.DeletedAt.Time
		dto.DeletedAt = &deletedAt
	}
	return dto
}

func ParseSyncHabitCompletion(completion *models.HabitCompletion) SyncHabitCompletionDTO {
	dto := SyncHabitCompletionDTO{}
	if parsed := ParseHabitCompletion(completion); parsed != nil {
		dto.HabitCompletionDTO = *parsed
	}
	if completion != nil && completion.DeletedAt.Valid {
		deletedAt := completion.DeletedAt.Time
		dto.DeletedAt = &deletedAt
	}
	return dto
}
