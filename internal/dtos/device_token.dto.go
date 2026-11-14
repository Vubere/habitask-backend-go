package dtos

import "habitask-backend-go/internal/models"

type DeviceTokenDTO struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

func ParseDeviceToken(token *models.DeviceToken) *DeviceTokenDTO {
	if token == nil {
		return nil
	}
	return &DeviceTokenDTO{
		ID:       token.ID,
		UserID:   token.UserID,
		Token:    token.Token,
		Platform: token.Platform,
	}
}

type DeviceTokenCreateDTO struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform"`
}
