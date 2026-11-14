package repositories

import (
	"time"

	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"

	"gorm.io/gorm"
)

type SyncRepository interface {
	GetChanges(userID string, since time.Time) (*dtos.SyncDTO, error)
}

type syncRepository struct {
	db *gorm.DB
}

func NewSyncRepository(db *gorm.DB) SyncRepository {
	return &syncRepository{db: db}
}

func (r *syncRepository) GetChanges(userID string, since time.Time) (*dtos.SyncDTO, error) {
	habits := []models.Habit{}
	habitQuery := r.db.Unscoped().Where("user_id = ?", userID)
	if !since.IsZero() {
		habitQuery = habitQuery.Where("updated_at >= ?", since)
	}
	if err := habitQuery.Find(&habits).Error; err != nil {
		return nil, err
	}

	tasks := []models.Task{}
	taskQuery := r.db.Unscoped().Where("user_id = ?", userID)
	if !since.IsZero() {
		taskQuery = taskQuery.Where("updated_at >= ?", since)
	}
	if err := taskQuery.Find(&tasks).Error; err != nil {
		return nil, err
	}

	completions := []models.HabitCompletion{}
	completionQuery := r.db.Unscoped().Where("user_id = ?", userID)
	if !since.IsZero() {
		completionQuery = completionQuery.Where("updated_at >= ?", since)
	}
	if err := completionQuery.Find(&completions).Error; err != nil {
		return nil, err
	}

	payload := &dtos.SyncDTO{
		Habits:           make([]dtos.SyncHabitDTO, 0, len(habits)),
		Tasks:            make([]dtos.SyncTaskDTO, 0, len(tasks)),
		HabitCompletions: make([]dtos.SyncHabitCompletionDTO, 0, len(completions)),
		ServerTime:       time.Now().UTC(),
	}
	for i := range habits {
		payload.Habits = append(payload.Habits, dtos.ParseSyncHabit(&habits[i]))
	}
	for i := range tasks {
		payload.Tasks = append(payload.Tasks, dtos.ParseSyncTask(&tasks[i]))
	}
	for i := range completions {
		payload.HabitCompletions = append(payload.HabitCompletions, dtos.ParseSyncHabitCompletion(&completions[i]))
	}
	return payload, nil
}
