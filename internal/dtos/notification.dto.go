package dtos

import (
	"habitask-backend-go/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type NotificationDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	IsRead      *bool     `json:"is_read"`
	Source      string    `json:"source"`
	SourceID    string    `json:"source_id"`
	UserID      string    `json:"user_id"`
	User        *UserDTO  `json:"user,omitempty"`
}

func ParseNotification(notification *models.Notification) *NotificationDTO {
	return &NotificationDTO{
		ID:          notification.ID,
		Title:       notification.Title,
		Description: notification.Description,
		Category:    notification.Category,
		Date:        notification.Date,
		IsRead:      notification.IsRead,
		Source:      notification.Source,
		SourceID:    notification.SourceID,
		UserID:      notification.UserID,
		User:        ParseUser(notification.User),
	}
}

type NotificationCreateDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Date        string `json:"date" binding:"required"`
	IsRead      *bool  `json:"is_read"`
	Source      string `json:"source" binding:"required"`
	SourceID    string `json:"source_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
}

func (n *NotificationCreateDTO) ToNotification(ctx *gin.Context) *models.Notification {
	notification := &models.Notification{
		Title:       n.Title,
		Description: n.Description,
		Category:    n.Category,
		IsRead:      n.IsRead,
		Source:      n.Source,
		SourceID:    n.SourceID,
	}
	if n.Date != "" {
		dateParsed, err := time.Parse("2006-01-02 15:04:05", n.Date)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid date format on field date"})
			return nil
		}
		notification.Date = dateParsed
	}
	return notification
}

type NotificationUpdateDTO struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Date        *string `json:"date"`
	IsRead      *bool   `json:"is_read"`
	Source      *string `json:"source"`
	SourceID    *string `json:"source_id"`
}

func (n *NotificationUpdateDTO) ToNotification(ctx *gin.Context) *models.Notification {
	notification := &models.Notification{}
	if n.Title != nil {
		notification.Title = *n.Title
	}
	if n.Description != nil {
		notification.Description = *n.Description
	}
	if n.Category != nil {
		notification.Category = *n.Category
	}
	if n.Date != nil {
		dateParsed, err := time.Parse("2006-01-02 15:04:05", *n.Date)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid date format on field date"})
			return nil
		}
		notification.Date = dateParsed
	}
	if n.IsRead != nil {
		notification.IsRead = n.IsRead
	}
	if n.Source != nil {
		notification.Source = *n.Source
	}
	if n.SourceID != nil {
		notification.SourceID = *n.SourceID
	}
	return notification
}
