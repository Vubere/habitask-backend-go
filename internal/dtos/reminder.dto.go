package dtos

import (
	"habitask-backend-go/internal/models"

	"github.com/gin-gonic/gin"
)

type ReminderDTO struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	Category            string   `json:"category"`
	Trigger             string   `json:"trigger"`
	TriggerValue        string   `json:"trigger_value"`
	TriggerIntervalUnit *string  `json:"trigger_interval_unit,omitempty"`
	IsRead              *bool    `json:"is_read"`
	Source              *string  `json:"source,omitempty"`
	SourceID            *string  `json:"source_id,omitempty"`
	SourceField         *string  `json:"source_field,omitempty"`
	SourceFieldType     *string  `json:"source_field_type,omitempty"`
	UserID              string   `json:"user_id"`
	User                *UserDTO `json:"user,omitempty"`
}

func ParseReminder(reminder *models.Reminder) *ReminderDTO {
	return &ReminderDTO{
		ID:                  reminder.ID,
		Title:               reminder.Title,
		Description:         reminder.Description,
		Category:            reminder.Category,
		Trigger:             reminder.Trigger,
		TriggerValue:        reminder.TriggerValue,
		TriggerIntervalUnit: reminder.TriggerIntervalUnit,
		IsRead:              reminder.IsRead,
		Source:              reminder.Source,
		SourceID:            reminder.SourceID,
		SourceField:         reminder.SourceField,
		SourceFieldType:     reminder.SourceFieldType,
		UserID:              reminder.UserID,
		User:                ParseUser(reminder.User),
	}
}

type ReminderCreateDTO struct {
	Title               string  `json:"title" binding:"required"`
	Description         string  `json:"description" binding:"required"`
	Category            string  `json:"category" binding:"required"`
	Trigger             string  `json:"trigger" binding:"required"`
	TriggerValue        string  `json:"trigger_value" binding:"required"`
	TriggerIntervalUnit *string `json:"trigger_interval_unit"`
	Source              *string `json:"source"`
	SourceID            *string `json:"source_id"`
	SourceField         *string `json:"source_field"`
	SourceFieldType     *string `json:"source_field_type"`
}

func (r *ReminderCreateDTO) ToReminder(ctx *gin.Context) *models.Reminder {
	reminder := &models.Reminder{
		Title:               r.Title,
		Description:         r.Description,
		Category:            r.Category,
		Trigger:             r.Trigger,
		TriggerValue:        r.TriggerValue,
		TriggerIntervalUnit: r.TriggerIntervalUnit,
		Source:              r.Source,
		SourceID:            r.SourceID,
		SourceField:         r.SourceField,
		SourceFieldType:     r.SourceFieldType,
	}
	if reminder.Trigger == "source" {
		if reminder.SourceID == nil || *reminder.SourceID == "" {
			ctx.JSON(400, gin.H{"error": "source_id is required for source trigger"})
			return nil
		}
		if reminder.SourceField == nil || *reminder.SourceField == "" {
			ctx.JSON(400, gin.H{"error": "source_field is required for source trigger"})
			return nil
		}
		if reminder.SourceFieldType == nil || *reminder.SourceFieldType == "" {
			ctx.JSON(400, gin.H{"error": "source_field_type is required for source trigger"})
			return nil
		}
	} else if reminder.Trigger == "interval" {
		if reminder.TriggerIntervalUnit == nil || *reminder.TriggerIntervalUnit == "" {
			ctx.JSON(400, gin.H{"error": "trigger_interval_unit is required for interval trigger"})
			return nil
		}
	}
	return reminder
}

type ReminderUpdateDTO struct {
	Title               *string `json:"title"`
	Description         *string `json:"description"`
	Category            *string `json:"category"`
	Trigger             *string `json:"trigger"`
	TriggerValue        *string `json:"trigger_value"`
	TriggerIntervalUnit *string `json:"trigger_interval_unit"`
	Source              *string `json:"source"`
	SourceID            *string `json:"source_id"`
	SourceField         *string `json:"source_field"`
	SourceFieldType     *string `json:"source_field_type"`
}

func (r *ReminderUpdateDTO) ToReminder() *models.Reminder {
	reminder := &models.Reminder{}
	if r.Title != nil {
		reminder.Title = *r.Title
	}
	if r.Description != nil {
		reminder.Description = *r.Description
	}
	if r.Category != nil {
		reminder.Category = *r.Category
	}
	if r.Trigger != nil {
		reminder.Trigger = *r.Trigger
	}
	if r.TriggerValue != nil {
		reminder.TriggerValue = *r.TriggerValue
	}
	if r.TriggerIntervalUnit != nil {
		reminder.TriggerIntervalUnit = r.TriggerIntervalUnit
	}
	if r.Source != nil {
		reminder.Source = r.Source
	}
	if r.SourceID != nil {
		reminder.SourceID = r.SourceID
	}
	if r.SourceField != nil {
		reminder.SourceField = r.SourceField
	}
	if r.SourceFieldType != nil {
		reminder.SourceFieldType = r.SourceFieldType
	}
	return reminder
}
