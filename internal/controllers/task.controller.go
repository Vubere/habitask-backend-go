package controllers

import (
	"errors"
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/lib/structs"

	"github.com/gin-gonic/gin"
)

type taskController struct {
	taskService services.TaskService
}

func NewTaskController(taskService services.TaskService) *taskController {
	return &taskController{taskService: taskService}
}

func (c *taskController) GetTasks(ctx *gin.Context) {
	pagination, query, err := c.getTaskPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	tasks, err := c.taskService.GetTasks(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": tasks, "pagination": pagination, "success": true, "message": "tasks retrieved successfully"})
}

func (c *taskController) GetTaskSummary(ctx *gin.Context) {
	pagination, query, err := c.getTaskPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	tasks, err := c.taskService.GetTaskSummary(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": tasks, "pagination": pagination, "success": true, "message": "tasks retrieved successfully"})
}

func (c *taskController) GetTaskById(ctx *gin.Context) {
	taskId := ctx.Param("id")
	task, err := c.taskService.GetTaskById(taskId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": task, "success": true, "message": "task retrieved successfully"})
}

func (c *taskController) CreateTask(ctx *gin.Context) {
	var taskInput dtos.TaskCreateDTO
	if err := ctx.ShouldBindJSON(&taskInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	task := taskInput.ToTask()
	err := c.taskService.CreateTask(task)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(201, gin.H{"message": "task created successfully", "success": true})
}

func (c *taskController) UpdateTask(ctx *gin.Context) {
	taskId := ctx.Param("id")
	var taskInput dtos.TaskUpdateDTO
	if err := ctx.ShouldBindJSON(&taskInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	taskToUpdate := taskInput.ToTask()
	taskToUpdate.ID = taskId
	err := c.taskService.UpdateTask(taskId, taskToUpdate)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "task updated successfully", "success": true})
}

func (c *taskController) DeleteTask(ctx *gin.Context) {
	taskId := ctx.Param("id")
	err := c.taskService.DeleteTask(taskId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "task deleted successfully", "success": true})
}

func (c *taskController) getTaskPaginationAndQuery(ctx *gin.Context) (*structs.PaginationAndSort, *models.TaskQuery, error) {
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
	query, ok := queryPassed.(*models.TaskQuery)
	if !ok {
		return nil, nil, errors.New("invalid query")
	}
	return pagination, query, nil
}
