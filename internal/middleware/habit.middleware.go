package middleware

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"
	"habitask-backend-go/pkg/lib/constants"

	"github.com/gin-gonic/gin"
)

func HabitMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		query := models.HabitQuery{
			Habit: models.Habit{
				UserID: helpers.ProcessUserIdFromContext(ctx),
				Name:   ctx.GetString("name"),
			},
		}
		trueVar := true
		falseVar := false
		if positive := ctx.Query("positive"); positive != "" {
			switch positive {
			case "true":
				query.IsPositive = &trueVar
			case "false":
				query.IsPositive = &falseVar
			}
		}
		if cat := ctx.Query("category"); cat != "" {
			query.Category = constants.Category(cat)
		}
		if ctx.Query("lasd_done_lte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("last_done_lte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.LastDoneLte = &cal
		}
		if ctx.Query("last_done_gte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("last_done_gte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.LastDoneGte = &cal
		}
		if ctx.Query("search") != "" {
			query.Search = ctx.Query("search")
		}
		if ctx.Query("group_by") != "" {
			query.GroupBy = ctx.Query("group_by")
		}
		if ctx.Query("date_group") != "" {
			query.DateGroup = ctx.Query("date_group")
		}

		paginationAndSort := helpers.GetPaginationAndSortFromContext(ctx)
		ctx.Set("paginationAndSort", &paginationAndSort)
		ctx.Set("query", &query)
		ctx.Next()
	}
}
