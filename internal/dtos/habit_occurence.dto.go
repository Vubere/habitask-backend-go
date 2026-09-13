package dtos

import (
	"errors"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"
)

type HabitOccurenceSummaryDTO struct {
	HabitID    string    `json:"habit_id"`
	Habit      *HabitDTO `json:"habit,omitempty"`
	Occurences int       `json:"occurences"`
	LastDone   *string   `json:"last_done"`
}

type HabitOccurenceCreateDTO struct {
	HabitID string `json:"habit_id" binding:"required"`
	Date    string `json:"date" binding:"required"`
	UserID  string `json:"user_id"`
}

func (h *HabitOccurenceCreateDTO) ToHabitOccurence() (*models.HabitOccurence, error) {
	habitOccurence := &models.HabitOccurence{
		HabitID: h.HabitID,
	}
	if h.UserID != "" {
		habitOccurence.UserID = h.UserID
	}
	if h.Date != "" {
		cal, err := helpers.ParseDate(h.Date)
		if err != nil {
			return nil, errors.New("invalid date format on field date")
		}
		habitOccurence.Date = cal
	}
	return habitOccurence, nil
}

type HabitOccurenceUpdateDTO struct {
	Date string `json:"date"`
}

func (h *HabitOccurenceUpdateDTO) ToHabitOccurence() *models.HabitOccurence {
	habitOccurence := &models.HabitOccurence{}
	if h.Date != "" {
		cal, err := helpers.ParseDate(h.Date)
		if err != nil {
			return nil
		}
		habitOccurence.Date = cal
	}
	return habitOccurence
}
