package controllers

import (
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

type DeviceTokenController interface {
	Register(ctx *gin.Context)
	Unregister(ctx *gin.Context)
}

type deviceTokenController struct {
	deviceTokenService services.DeviceTokenService
}

func NewDeviceTokenController(deviceTokenService services.DeviceTokenService) *deviceTokenController {
	return &deviceTokenController{deviceTokenService: deviceTokenService}
}

func (c *deviceTokenController) Register(ctx *gin.Context) {
	var input dtos.DeviceTokenCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	token := &models.DeviceToken{Token: input.Token, Platform: input.Platform}
	if err := c.deviceTokenService.Register(ctx.GetString("userId"), token); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(201, gin.H{"message": "device registered successfully", "success": true})
}

func (c *deviceTokenController) Unregister(ctx *gin.Context) {
	if err := c.deviceTokenService.Unregister(ctx.Param("token")); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "device unregistered successfully", "success": true})
}
