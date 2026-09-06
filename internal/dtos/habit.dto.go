package dtos

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/lib/constants"
	"time"

	"github.com/gin-gonic/gin"
)

type HabitDTO struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Category    constants.Category `json:"category"`
	Pros        string             `json:"pros"`
	Cons        string             `json:"cons"`
	IsPositive  *bool              `json:"is_positive"`
	LastDone    *time.Time         `json:"last_done"`
	User        *UserDTO           `json:"user,omitempty"`
}

func ParseHabit(habit *models.Habit) *HabitDTO {
	return &HabitDTO{
		ID:          habit.ID,
		Name:        habit.Name,
		Description: habit.Description,
		Category:    habit.Category,
		Pros:        habit.Pros,
		Cons:        habit.Cons,
		IsPositive:  habit.IsPositive,
		LastDone:    habit.LastDone,
		User:        ParseUser(habit.User),
	}
}

type HabitCreateDTO struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description" binding:"required"`
	Category    constants.Category `json:"category" binding:"required"`
	Pros        string             `json:"pros" binding:"required"`
	Cons        string             `json:"cons" binding:"required"`
	IsPositive  *bool              `json:"is_positive"`
	LastDone    *string            `json:"last_done"`
}

func (h *HabitCreateDTO) ToHabit(ctx *gin.Context) *models.Habit {
	habit := &models.Habit{
		Name:        h.Name,
		Description: h.Description,
		Category:    h.Category,
		Pros:        h.Pros,
		Cons:        h.Cons,
		IsPositive:  h.IsPositive,
		LastDone:    nil,
	}
	if h.LastDone != nil {
		ld, err := time.Parse("2006-01-02 15:04:05", *h.LastDone)
		if err != nil {
			if ctx != nil {
				ctx.JSON(400, gin.H{"error": "invalid date format on field last_done"})
			}
			return nil
		}
		habit.LastDone = &ld
	}
	return habit
}

type HabitUpdateDTO struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Category    *constants.Category `json:"category"`
	Pros        *string             `json:"pros"`
	Cons        *string             `json:"cons"`
	IsPositive  *bool               `json:"is_positive"`
	LastDone    *string             `json:"last_done"`
}

func (h *HabitUpdateDTO) ToHabit() *models.Habit {
	habit := &models.Habit{}
	if h.Name != nil {
		habit.Name = *h.Name
	}
	if h.Description != nil {
		habit.Description = *h.Description
	}
	if h.Category != nil {
		habit.Category = *h.Category
	}
	if h.Pros != nil {
		habit.Pros = *h.Pros
	}
	if h.Cons != nil {
		habit.Cons = *h.Cons
	}
	if h.IsPositive != nil {
		habit.IsPositive = h.IsPositive
	}
	if h.LastDone != nil {
		ld, err := time.Parse("2006-01-02 15:04:05", *h.LastDone)
		if err != nil {
			return nil
		}
		habit.LastDone = &ld
	}
	return habit
}
