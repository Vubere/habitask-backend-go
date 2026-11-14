package middleware

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"

	"github.com/gin-gonic/gin"
)

func HabitCompletionMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		query := models.HabitCompletionQuery{
			HabitCompletion: models.HabitCompletion{
				UserID: helpers.ProcessUserIdFromContext(ctx),
			},
		}
		if habitId := ctx.Query("habit_id"); habitId != "" {
			query.HabitID = habitId
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
		if ctx.Query("group_by") != "" {
			query.GroupBy = ctx.Query("group_by")
		}
		if ctx.Query("date_group") != "" {
			query.DateGroup = ctx.Query("date_group")
		}

		paginationAndSort := helpers.GetPaginationAndSortFromContext(ctx)
		ctx.Set("pagination", &paginationAndSort)
		ctx.Set("query", &query)
		ctx.Next()
	}
}
