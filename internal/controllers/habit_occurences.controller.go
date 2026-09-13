package controllers

import (
	"errors"
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/lib/structs"

	"github.com/gin-gonic/gin"
)

type HabitOccurencesController interface {
	GetHabitOccurences(ctx *gin.Context)
	GetHabitOccurenceSummary(ctx *gin.Context)
	GetHabitOccurenceById(ctx *gin.Context)
	CreateHabitOccurence(ctx *gin.Context)
	UpdateHabitOccurence(ctx *gin.Context)
	DeleteHabitOccurence(ctx *gin.Context)
}

type habitOccurencesController struct {
	habitOccurenceService services.HabitOccurenceService
}

func NewHabitOccurencesController(habitOccurenceService services.HabitOccurenceService) *habitOccurencesController {
	return &habitOccurencesController{habitOccurenceService: habitOccurenceService}
}

func (c *habitOccurencesController) GetHabitOccurences(ctx *gin.Context) {
	pagination, query, err := c.getHabitOccurencePaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitOccurences, err := c.habitOccurenceService.GetHabitOccurences(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": habitOccurences, "pagination": pagination, "success": true, "message": "habit occurences retrieved successfully"})
}

func (c *habitOccurencesController) GetHabitOccurenceSummary(ctx *gin.Context) {
	pagination, query, err := c.getHabitOccurencePaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitOccurenceSummaries, err := c.habitOccurenceService.GetHabitOccurenceSummary(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": habitOccurenceSummaries, "pagination": pagination, "success": true, "message": "habit occurence summaries retrieved successfully"})
}

func (c *habitOccurencesController) GetHabitOccurenceById(ctx *gin.Context) {
	habitOccurenceId := ctx.Param("id")
	habitOccurence, err := c.habitOccurenceService.GetHabitOccurenceById(habitOccurenceId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": habitOccurence, "success": true, "message": "habit occurence retrieved successfully"})
}

func (c *habitOccurencesController) CreateHabitOccurence(ctx *gin.Context) {
	var habitOccurenceInput dtos.HabitOccurenceCreateDTO
	if err := ctx.ShouldBindJSON(&habitOccurenceInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitOccurence, err := habitOccurenceInput.ToHabitOccurence()
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitOccurence.UserID = ctx.GetString("userId")

	err = c.habitOccurenceService.CreateHabitOccurence(habitOccurence)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(201, gin.H{"message": "habit occurence created successfully", "success": true})
}

func (c *habitOccurencesController) UpdateHabitOccurence(ctx *gin.Context) {
	habitOccurenceId := ctx.Param("id")
	var habitOccurenceInput dtos.HabitOccurenceUpdateDTO
	if err := ctx.ShouldBindJSON(&habitOccurenceInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitOccurenceToUpdate := habitOccurenceInput.ToHabitOccurence()

	habitOccurenceToUpdate.ID = habitOccurenceId
	err := c.habitOccurenceService.UpdateHabitOccurence(habitOccurenceId, habitOccurenceToUpdate)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "habit occurence updated successfully", "success": true})
}

func (c *habitOccurencesController) DeleteHabitOccurence(ctx *gin.Context) {
	habitOccurenceId := ctx.Param("id")
	err := c.habitOccurenceService.DeleteHabitOccurence(habitOccurenceId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "habit occurence deleted successfully", "success": true})
}

func (c *habitOccurencesController) getHabitOccurencePaginationAndQuery(ctx *gin.Context) (*structs.PaginationAndSort, *models.HabitOccurenceQuery, error) {
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
	query, ok := queryPassed.(*models.HabitOccurenceQuery)
	if !ok {
		return nil, nil, errors.New("invalid query")
	}
	return pagination, query, nil
}
