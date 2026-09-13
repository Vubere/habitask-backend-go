package controllers

import (
	"errors"
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/lib/structs"

	"github.com/gin-gonic/gin"
)

type HabitActionController interface {
	GetHabitActions(ctx *gin.Context)
	GetHabitActionSummary(ctx *gin.Context)
	GetHabitActionById(ctx *gin.Context)
	CreateHabitAction(ctx *gin.Context)
	UpdateHabitAction(ctx *gin.Context)
	DeleteHabitAction(ctx *gin.Context)
}

type habitActionController struct {
	habitActionService services.HabitActionService
}

func NewHabitActionController(habitActionService services.HabitActionService) *habitActionController {
	return &habitActionController{habitActionService: habitActionService}
}

func (c *habitActionController) GetHabitActions(ctx *gin.Context) {
	pagination, query, err := c.getHabitActionPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitActions, err := c.habitActionService.GetHabitActions(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": habitActions, "pagination": pagination, "success": true, "message": "habit actions retrieved successfully"})
}

func (c *habitActionController) GetHabitActionSummary(ctx *gin.Context) {
	pagination, query, err := c.getHabitActionPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitActionSummaries, err := c.habitActionService.GetHabitActionSummary(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": habitActionSummaries, "pagination": pagination, "success": true, "message": "habit action summaries retrieved successfully"})
}

func (c *habitActionController) GetHabitActionById(ctx *gin.Context) {
	habitActionId := ctx.Param("id")
	habitAction, err := c.habitActionService.GetHabitActionById(habitActionId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": habitAction, "success": true, "message": "habit action retrieved successfully"})
}

func (c *habitActionController) CreateHabitAction(ctx *gin.Context) {
	var habitActionInput dtos.HabitActionCreateDTO
	if err := ctx.ShouldBindJSON(&habitActionInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitAction := habitActionInput.ToHabitAction()

	err := c.habitActionService.CreateHabitAction(habitAction)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(201, gin.H{"message": "habit action created successfully", "success": true})
}

func (c *habitActionController) UpdateHabitAction(ctx *gin.Context) {
	habitActionId := ctx.Param("id")
	var habitActionInput dtos.HabitActionUpdateDTO
	if err := ctx.ShouldBindJSON(&habitActionInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitActionToUpdate := habitActionInput.ToHabitAction()

	habitActionToUpdate.ID = habitActionId
	err := c.habitActionService.UpdateHabitAction(habitActionId, habitActionToUpdate)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "habit action updated successfully", "success": true})
}

func (c *habitActionController) DeleteHabitAction(ctx *gin.Context) {
	habitActionId := ctx.Param("id")
	err := c.habitActionService.DeleteHabitAction(habitActionId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "habit action deleted successfully", "success": true})
}

func (c *habitActionController) getHabitActionPaginationAndQuery(ctx *gin.Context) (*structs.PaginationAndSort, *models.HabitActionQuery, error) {
	paginationAndSortPassed, exists := ctx.Get("pagination")
	if !exists {
		return nil, nil, errors.New("invalid pagination")
	}
	pagination, ok := paginationAndSortPassed.(*structs.PaginationAndSort)
	if !ok {
		return nil, nil, errors.New("invalid pagination")
	}
	queryPassed, exists := ctx.Get("query")
	if !exists {
		return nil, nil, errors.New("invalid query")
	}
	query, ok := queryPassed.(*models.HabitActionQuery)
	if !ok {
		return nil, nil, errors.New("invalid query")
	}
	return pagination, query, nil
}
