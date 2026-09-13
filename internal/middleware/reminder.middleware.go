package middleware

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"

	"github.com/gin-gonic/gin"
)

func ReminderMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		query := models.ReminderQuery{
			Reminder: models.Reminder{
				UserID: helpers.ProcessUserIdFromContext(ctx),
			},
		}
		if cat := ctx.Query("category"); cat != "" {
			query.Category = cat
		}
		if pur := ctx.Query("purpose"); pur != "" {
			query.Purpose = pur
		}
		if trig := ctx.Query("trigger"); trig != "" {
			query.Trigger = trig
		}
		if trigVal := ctx.Query("trigger_value"); trigVal != "" {
			query.TriggerValue = trigVal
		}
		if trigUnit := ctx.Query("trigger_interval_unit"); trigUnit != "" {
			query.TriggerIntervalUnit = &trigUnit
		}
		if src := ctx.Query("source"); src != "" {
			query.Source = &src
		}
		if srcID := ctx.Query("source_id"); srcID != "" {
			query.SourceID = &srcID
		}
		if srcField := ctx.Query("source_field"); srcField != "" {
			query.SourceField = &srcField
		}
		if srcFieldType := ctx.Query("source_field_type"); srcFieldType != "" {
			query.SourceFieldType = &srcFieldType
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
