package controllers

import (
	"errors"
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/lib/structs"

	"github.com/gin-gonic/gin"
)

type ReminderController interface {
	GetReminders(ctx *gin.Context)
	GetReminderById(ctx *gin.Context)
	CreateReminder(ctx *gin.Context)
	UpdateReminder(ctx *gin.Context)
	DeleteReminder(ctx *gin.Context)
}

type reminderController struct {
	reminderService services.ReminderService
}

func NewReminderController(reminderService services.ReminderService) *reminderController {
	return &reminderController{reminderService: reminderService}
}

func (c *reminderController) GetReminders(ctx *gin.Context) {
	pagination, query, err := c.getReminderPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	reminders, err := c.reminderService.GetReminders(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": reminders, "pagination": pagination, "success": true, "message": "reminders retrieved successfully"})
}

func (c *reminderController) GetReminderById(ctx *gin.Context) {
	reminderId := ctx.Param("id")
	reminder, err := c.reminderService.GetReminderById(reminderId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": reminder, "success": true, "message": "reminder retrieved successfully"})
}

func (c *reminderController) CreateReminder(ctx *gin.Context) {
	var reminderInput dtos.ReminderCreateDTO
	if err := ctx.ShouldBindJSON(&reminderInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	reminder, err := reminderInput.ToReminder()
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	err = c.reminderService.CreateReminder(reminder)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(201, gin.H{"message": "reminder created successfully", "success": true})
}

func (c *reminderController) UpdateReminder(ctx *gin.Context) {
	reminderId := ctx.Param("id")
	var reminderInput dtos.ReminderUpdateDTO
	if err := ctx.ShouldBindJSON(&reminderInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	reminderToUpdate := reminderInput.ToReminder()

	reminderToUpdate.ID = reminderId
	err := c.reminderService.UpdateReminder(reminderId, reminderToUpdate)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "reminder updated successfully", "success": true})
}

func (c *reminderController) DeleteReminder(ctx *gin.Context) {
	reminderId := ctx.Param("id")
	err := c.reminderService.DeleteReminder(reminderId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "reminder deleted successfully", "success": true})
}

func (c *reminderController) getReminderPaginationAndQuery(ctx *gin.Context) (*structs.PaginationAndSort, *models.ReminderQuery, error) {
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
	query, ok := queryPassed.(*models.ReminderQuery)
	if !ok {
		return nil, nil, errors.New("invalid query")
	}
	return pagination, query, nil
}
