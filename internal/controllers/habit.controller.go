package controllers

import (
	"errors"
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/lib/structs"

	"github.com/gin-gonic/gin"
)

type HabitController interface {
	GetHabits(ctx *gin.Context)
	GetHabitSummary(ctx *gin.Context)
	GetHabitById(ctx *gin.Context)
	CreateHabit(ctx *gin.Context)
	UpdateHabit(ctx *gin.Context)
	DeleteHabit(ctx *gin.Context)
}

type habitController struct {
	habitService services.HabitService
}

func NewHabitController(habitService services.HabitService) *habitController {
	return &habitController{habitService: habitService}
}

func (c *habitController) GetHabits(ctx *gin.Context) {
	pagination, query, err := c.getHabitPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habits, err := c.habitService.GetHabits(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": habits, "pagination": pagination, "success": true, "message": "habits retrieved successfully"})
}

func (c *habitController) GetHabitSummary(ctx *gin.Context) {
	pagination, query, err := c.getHabitPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitSummaries, err := c.habitService.GetHabitSummary(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": habitSummaries, "pagination": pagination, "success": true, "message": "habit summaries retrieved successfully"})
}

func (c *habitController) GetHabitById(ctx *gin.Context) {
	habitId := ctx.Param("id")
	habit, err := c.habitService.GetHabitById(habitId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": habit, "success": true, "message": "habit retrieved successfully"})
}

func (c *habitController) CreateHabit(ctx *gin.Context) {
	var habitInput dtos.HabitCreateDTO
	if err := ctx.ShouldBindJSON(&habitInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habit, err := habitInput.ToHabit()
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	err = c.habitService.CreateHabit(habit)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(201, gin.H{"message": "habit created successfully", "success": true})
}

func (c *habitController) UpdateHabit(ctx *gin.Context) {
	habitId := ctx.Param("id")
	var habitInput dtos.HabitUpdateDTO
	if err := ctx.ShouldBindJSON(&habitInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	habitToUpdate := habitInput.ToHabit()
	habitToUpdate.ID = habitId
	err := c.habitService.UpdateHabit(habitId, habitToUpdate)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "habit updated successfully", "success": true})
}

func (c *habitController) DeleteHabit(ctx *gin.Context) {
	habitId := ctx.Param("id")
	err := c.habitService.DeleteHabit(habitId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "habit deleted successfully", "success": true})
}

func (c *habitController) getHabitPaginationAndQuery(ctx *gin.Context) (*structs.PaginationAndSort, *models.HabitQuery, error) {
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
	query, ok := queryPassed.(*models.HabitQuery)
	if !ok {
		return nil, nil, errors.New("invalid query")
	}
	return pagination, query, nil
}
