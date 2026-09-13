package middleware

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"
	"habitask-backend-go/pkg/lib/constants"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TaskMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		query := models.TaskQuery{
			Task: models.Task{
				Name:   ctx.GetString("name"),
				UserID: helpers.ProcessUserIdFromContext(ctx),
			},
		}
		trueVar := true
		falseVar := false
		if cat := ctx.Query("category"); cat != "" {
			query.Category = constants.Category(cat)
		}
		if pr := ctx.Query("priority"); pr != "" {
			query.Priority = pr
		}
		if completed := ctx.Query("completed"); completed != "" {
			switch completed {
			case "true":
				query.IsCompleted = &trueVar
			case "false":
				query.IsCompleted = &falseVar
			}
		}
		if ctx.Query("completed_at_lte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("completed_at_lte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.CompletedAtLte = &cal
		}
		if ctx.Query("completed_at_gte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("completed_at_gte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.CompletedAtGte = &cal
		}
		if ctx.Query("due_at_lte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("due_at_lte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.DueAtLte = &cal
		}
		if ctx.Query("due_at_gte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("due_at_gte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.DueAtGte = &cal
		}
		if edl := ctx.Query("estimated_duration_lte"); edl != "" {
			intEdl, err := strconv.Atoi(edl)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.EstimatedDurationLte = &intEdl
		}
		if edg := ctx.Query("estimated_duration_gte"); edg != "" {
			intEdg, err := strconv.Atoi(edg)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.EstimatedDurationGte = &intEdg
		}
		if adl := ctx.Query("actual_duration_lte"); adl != "" {
			intAdl, err := strconv.Atoi(adl)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.ActualDurationLte = &intAdl
		}
		if adg := ctx.Query("actual_duration_gte"); adg != "" {
			intAdg, err := strconv.Atoi(adg)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.ActualDurationGte = &intAdg
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
