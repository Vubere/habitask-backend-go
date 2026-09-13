package middleware

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"

	"github.com/gin-gonic/gin"
)

func NotificationMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		query := models.NotificationQuery{
			Notification: models.Notification{
				UserID: helpers.ProcessUserIdFromContext(ctx),
			},
		}
		if cat := ctx.Query("category"); cat != "" {
			query.Category = cat
		}
		if src := ctx.Query("source"); src != "" {
			query.Source = src
		}
		if srcID := ctx.Query("source_id"); srcID != "" {
			query.SourceID = srcID
		}
		if ctx.Query("date_lte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("date_lte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.DateLte = &cal
		}
		if ctx.Query("date_gte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("date_gte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.DateGte = &cal
		}
		if ctx.Query("search") != "" {
			query.Search = ctx.Query("search")
		}

		paginationAndSort := helpers.GetPaginationAndSortFromContext(ctx)
		ctx.Set("paginationAndSort", &paginationAndSort)
		ctx.Set("query", &query)
		ctx.Next()
	}
}
