package controllers

import (
	"time"

	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/helpers"

	"github.com/gin-gonic/gin"
)

type SyncController interface {
	GetChanges(ctx *gin.Context)
}

type syncController struct {
	syncService services.SyncService
}

func NewSyncController(syncService services.SyncService) *syncController {
	return &syncController{syncService: syncService}
}

func (c *syncController) GetChanges(ctx *gin.Context) {
	since := time.Time{}
	if raw := ctx.Query("since"); raw != "" {
		parsed, err := helpers.ParseDate(raw)
		if err != nil {
			ctx.JSON(400, gin.H{"message": "invalid since date format", "success": false})
			return
		}
		since = parsed
	}
	payload, err := c.syncService.GetChanges(ctx.GetString("userId"), since)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": payload, "success": true, "message": "changes retrieved successfully"})
}
