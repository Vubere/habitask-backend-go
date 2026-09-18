package controllers

import (
	"errors"
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/lib/structs"

	"github.com/gin-gonic/gin"
)

type HabitTaskController interface {
	GetHabitTaskDisplayDetails(ctx *gin.Context)
}

type habitTaskController struct {
	habitTaskService services.HabitTaskService
}

func NewHabitTaskController(habitTaskService services.HabitTaskService) *habitTaskController {
	return &habitTaskController{habitTaskService: habitTaskService}
}

func (c *habitTaskController) GetHabitTaskDisplayDetails(ctx *gin.Context) {
	pagination, query, err := c.getHabitTaskPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	query.UserID = ctx.GetString("userId")
	habitTasks, err := c.habitTaskService.GetHabitTaskDisplayDetails(*query, pagination, ctx.Query("search"))
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": habitTasks, "pagination": pagination, "success": true, "message": "habit tasks retrieved successfully"})
}

func (c *habitTaskController) getHabitTaskPaginationAndQuery(ctx *gin.Context) (*structs.PaginationAndSort, *dtos.HabitTaskFilterDTO, error) {
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
	query, ok := queryPassed.(*dtos.HabitTaskFilterDTO)
	if !ok {
		return nil, nil, errors.New("invalid query")
	}
	return pagination, query, nil
}
