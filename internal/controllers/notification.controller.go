package controllers

import (
	"errors"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/lib/structs"

	"github.com/gin-gonic/gin"
)

type NotificationController interface {
	GetNotifications(ctx *gin.Context)
	GetNotificationSummary(ctx *gin.Context)
	GetNotificationById(ctx *gin.Context)
	MarkNotificationAsRead(ctx *gin.Context)
	MarkAllNotificationsAsRead(ctx *gin.Context)
}

type notificationController struct {
	notificationService services.NotificationService
}

func NewNotificationController(notificationService services.NotificationService) *notificationController {
	return &notificationController{notificationService: notificationService}
}

func (c *notificationController) GetNotifications(ctx *gin.Context) {
	pagination, query, err := c.getNotificationPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	notifications, err := c.notificationService.GetNotifications(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": notifications, "pagination": pagination, "success": true, "message": "notifications retrieved successfully"})
}

func (c *notificationController) GetNotificationSummary(ctx *gin.Context) {
	pagination, query, err := c.getNotificationPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	notificationSummaries, err := c.notificationService.GetNotificationSummary(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": notificationSummaries, "pagination": pagination, "success": true, "message": "notification summaries retrieved successfully"})
}

func (c *notificationController) GetNotificationById(ctx *gin.Context) {
	notificationId := ctx.Param("id")
	notification, err := c.notificationService.GetNotificationById(notificationId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": notification, "success": true, "message": "notification retrieved successfully"})
}

func (c *notificationController) MarkNotificationAsRead(ctx *gin.Context) {
	notificationId := ctx.Param("id")
	err := c.notificationService.MarkNotificationAsRead(notificationId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "notification marked as read successfully", "success": true})
}

func (c *notificationController) MarkAllNotificationsAsRead(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	err := c.notificationService.MarkAllNotificationsAsRead(userId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "all notifications marked as read successfully", "success": true})
}

func (c *notificationController) getNotificationPaginationAndQuery(ctx *gin.Context) (*structs.PaginationAndSort, *models.NotificationQuery, error) {
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
	query, ok := queryPassed.(*models.NotificationQuery)
	if !ok {
		return nil, nil, errors.New("invalid query")
	}
	return pagination, query, nil
}
