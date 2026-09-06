package dtos

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/lib/constants"

	"github.com/gin-gonic/gin"
)

type HabitActionDTO struct {
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	Category     constants.Category `json:"category"`
	CostIncurred float64            `json:"cost_incurred"`
	IsPositive   *bool              `json:"is_positive"`
	User         *UserDTO           `json:"user"`
}

func ParseHabitAction(habitAction *models.HabitAction) *HabitActionDTO {
	return &HabitActionDTO{
		ID:           habitAction.ID,
		Title:        habitAction.Title,
		Description:  habitAction.Description,
		Category:     habitAction.Category,
		CostIncurred: habitAction.CostIncurred,
		IsPositive:   habitAction.IsPositive,
		User:         ParseUser(habitAction.User),
	}
}

type HabitActionCreateDTO struct {
	Title        string             `json:"title" binding:"required"`
	Description  string             `json:"description" binding:"required"`
	Category     constants.Category `json:"category" binding:"required"`
	CostIncurred float64            `json:"cost_incurred" binding:"required"`
	IsPositive   *bool              `json:"is_positive"`
}

func (h *HabitActionCreateDTO) ToHabitAction(ctx *gin.Context) *models.HabitAction {
	habitAction := &models.HabitAction{
		Title:        h.Title,
		Description:  h.Description,
		Category:     h.Category,
		CostIncurred: h.CostIncurred,
		IsPositive:   h.IsPositive,
	}
	return habitAction
}

type HabitActionUpdateDTO struct {
	Title        *string             `json:"title"`
	Description  *string             `json:"description"`
	Category     *constants.Category `json:"category"`
	CostIncurred *float64            `json:"cost_incurred"`
	IsPositive   *bool               `json:"is_positive"`
}

func (h *HabitActionUpdateDTO) ToHabitAction() *models.HabitAction {
	habitAction := &models.HabitAction{}
	if h.Title != nil {
		habitAction.Title = *h.Title
	}
	if h.Description != nil {
		habitAction.Description = *h.Description
	}
	if h.Category != nil {
		habitAction.Category = *h.Category
	}
	if h.CostIncurred != nil {
		habitAction.CostIncurred = *h.CostIncurred
	}
	if h.IsPositive != nil {
		habitAction.IsPositive = h.IsPositive
	}
	return habitAction
}
