package controllers

import (
	"errors"
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/lib/structs"

	"github.com/gin-gonic/gin"
)

type TaskStepsController interface {
	GetTaskSteps(ctx *gin.Context)
	GetTaskStepSummary(ctx *gin.Context)
	GetTaskStepById(ctx *gin.Context)
	CreateTaskStep(ctx *gin.Context)
	UpdateTaskStep(ctx *gin.Context)
	DeleteTaskStep(ctx *gin.Context)
}

type taskStepsController struct {
	taskStepsService services.TaskStepService
}

func NewTaskStepsController(taskStepsService services.TaskStepService) *taskStepsController {
	return &taskStepsController{taskStepsService: taskStepsService}
}

func (c *taskStepsController) GetTaskSteps(ctx *gin.Context) {
	pagination, query, err := c.getTaskStepsPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	taskSteps, err := c.taskStepsService.GetTaskSteps(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": taskSteps, "pagination": pagination, "success": true, "message": "task steps retrieved successfully"})
}

func (c *taskStepsController) GetTaskStepSummary(ctx *gin.Context) {
	pagination, query, err := c.getTaskStepsPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	taskStepSummaries, err := c.taskStepsService.GetTaskStepSummary(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": taskStepSummaries, "pagination": pagination, "success": true, "message": "task step summaries retrieved successfully"})
}

func (c *taskStepsController) GetTaskStepById(ctx *gin.Context) {
	taskStepId := ctx.Param("id")
	taskStep, err := c.taskStepsService.GetTaskStepById(taskStepId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": taskStep, "success": true, "message": "task step retrieved successfully"})
}

func (c *taskStepsController) CreateTaskStep(ctx *gin.Context) {
	var taskStepInput dtos.TaskStepCreateDTO
	if err := ctx.ShouldBindJSON(&taskStepInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	taskStep, err := taskStepInput.ToTaskStep()
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	err = c.taskStepsService.CreateTaskStep(taskStep)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(201, gin.H{"message": "task step created successfully", "success": true})
}

func (c *taskStepsController) UpdateTaskStep(ctx *gin.Context) {
	taskStepId := ctx.Param("id")
	var taskStepInput dtos.TaskStepUpdateDTO
	if err := ctx.ShouldBindJSON(&taskStepInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	taskStepToUpdate, err := taskStepInput.ToTaskStep()
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	taskStepToUpdate.ID = taskStepId
	err = c.taskStepsService.UpdateTaskStep(taskStepId, taskStepToUpdate)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "task step updated successfully", "success": true})
}

func (c *taskStepsController) DeleteTaskStep(ctx *gin.Context) {
	taskStepId := ctx.Param("id")
	err := c.taskStepsService.DeleteTaskStep(taskStepId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "task step deleted successfully", "success": true})
}

func (c *taskStepsController) getTaskStepsPaginationAndQuery(ctx *gin.Context) (*structs.PaginationAndSort, *models.TaskStepsQuery, error) {
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
	query, ok := queryPassed.(*models.TaskStepsQuery)
	if !ok {
		return nil, nil, errors.New("invalid query")
	}
	return pagination, query, nil
}
