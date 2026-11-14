package controllers

import (
	"errors"
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/lib/structs"

	"github.com/gin-gonic/gin"
)

type HabitCompletionController interface {
	GetHabitCompletions(ctx *gin.Context)
	GetHabitCompletionSummary(ctx *gin.Context)
	GetHabitCompletionById(ctx *gin.Context)
	CreateHabitCompletion(ctx *gin.Context)
	LogCompletion(ctx *gin.Context)
	UpdateHabitCompletion(ctx *gin.Context)
	DeleteHabitCompletion(ctx *gin.Context)
}

type habitCompletionController struct {
	habitCompletionService services.HabitCompletionService
}

func NewHabitCompletionController(habitCompletionService services.HabitCompletionService) *habitCompletionController {
	return &habitCompletionController{habitCompletionService: habitCompletionService}
}

func (c *habitCompletionController) GetHabitCompletions(ctx *gin.Context) {
	pagination, query, err := c.getHabitCompletionPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	completions, err := c.habitCompletionService.GetHabitCompletions(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": dtos.ParseHabitCompletions(completions), "pagination": pagination, "success": true, "message": "habit completions retrieved successfully"})
}

func (c *habitCompletionController) GetHabitCompletionSummary(ctx *gin.Context) {
	pagination, query, err := c.getHabitCompletionPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	summaries, err := c.habitCompletionService.GetHabitCompletionSummary(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": summaries, "pagination": pagination, "success": true, "message": "habit completion summaries retrieved successfully"})
}

func (c *habitCompletionController) GetHabitCompletionById(ctx *gin.Context) {
	completionId := ctx.Param("id")
	completion, err := c.habitCompletionService.GetHabitCompletionById(completionId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": dtos.ParseHabitCompletion(completion), "success": true, "message": "habit completion retrieved successfully"})
}

func (c *habitCompletionController) CreateHabitCompletion(ctx *gin.Context) {
	var input dtos.HabitCompletionCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	completion, err := input.ToHabitCompletion()
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	completion.UserID = ctx.GetString("userId")
	if err := c.habitCompletionService.CreateHabitCompletion(completion); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(201, gin.H{"message": "habit completion created successfully", "success": true})
}

func (c *habitCompletionController) LogCompletion(ctx *gin.Context) {
	habitId := ctx.Param("id")
	var input dtos.HabitCompletionCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	completion, err := input.ToHabitCompletion()
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	completion.HabitID = habitId
	completion.UserID = ctx.GetString("userId")
	if err := c.habitCompletionService.LogHabitCompletion(completion); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "habit completion logged successfully", "success": true})
}

func (c *habitCompletionController) UpdateHabitCompletion(ctx *gin.Context) {
	completionId := ctx.Param("id")
	var input dtos.HabitCompletionUpdateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	completion, err := input.ToHabitCompletion()
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	completion.ID = completionId
	if err := c.habitCompletionService.UpdateHabitCompletion(completionId, completion); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "habit completion updated successfully", "success": true})
}

func (c *habitCompletionController) DeleteHabitCompletion(ctx *gin.Context) {
	completionId := ctx.Param("id")
	if err := c.habitCompletionService.DeleteHabitCompletion(completionId); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "habit completion deleted successfully", "success": true})
}

func (c *habitCompletionController) getHabitCompletionPaginationAndQuery(ctx *gin.Context) (*structs.PaginationAndSort, *models.HabitCompletionQuery, error) {
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
	query, ok := queryPassed.(*models.HabitCompletionQuery)
	if !ok {
		return nil, nil, errors.New("invalid query")
	}
	return pagination, query, nil
}
