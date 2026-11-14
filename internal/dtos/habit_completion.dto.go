package dtos

import (
	"errors"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"
	"time"
)

type HabitCompletionDTO struct {
	ID          string    `json:"id"`
	HabitID     string    `json:"habit_id"`
	UserID      string    `json:"user_id"`
	Date        string    `json:"date"`
	Count       int       `json:"count"`
	CompletedAt time.Time `json:"completed_at"`
}

func ParseHabitCompletion(completion *models.HabitCompletion) *HabitCompletionDTO {
	if completion == nil {
		return nil
	}
	dto := &HabitCompletionDTO{
		ID:          completion.ID,
		HabitID:     completion.HabitID,
		UserID:      completion.UserID,
		Count:       completion.Count,
		CompletedAt: completion.CompletedAt,
	}
	if !completion.Date.IsZero() {
		dto.Date = completion.Date.Format("2006-01-02")
	}
	return dto
}

func ParseHabitCompletions(completions []models.HabitCompletion) []HabitCompletionDTO {
	result := make([]HabitCompletionDTO, 0, len(completions))
	for i := range completions {
		result = append(result, *ParseHabitCompletion(&completions[i]))
	}
	return result
}

type HabitCompletionCreateDTO struct {
	HabitID string `json:"habit_id"`
	Date    string `json:"date"`
	Count   int    `json:"count"`
}

func (h *HabitCompletionCreateDTO) ToHabitCompletion() (*models.HabitCompletion, error) {
	completion := &models.HabitCompletion{
		HabitID:     h.HabitID,
		Count:       h.Count,
		CompletedAt: time.Now(),
	}
	if h.Date != "" {
		date, err := helpers.ParseDate(h.Date)
		if err != nil {
			return nil, errors.New("invalid date format on field date")
		}
		completion.Date = date
	} else {
		completion.Date = time.Now()
	}
	return completion, nil
}

type HabitCompletionUpdateDTO struct {
	Date  *string `json:"date"`
	Count *int    `json:"count"`
}

func (h *HabitCompletionUpdateDTO) ToHabitCompletion() (*models.HabitCompletion, error) {
	completion := &models.HabitCompletion{}
	if h.Date != nil {
		date, err := helpers.ParseDate(*h.Date)
		if err != nil {
			return nil, errors.New("invalid date format on field date")
		}
		completion.Date = date
	}
	if h.Count != nil {
		completion.Count = *h.Count
	}
	completion.CompletedAt = time.Now()
	return completion, nil
}
