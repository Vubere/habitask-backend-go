package middleware

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TaskStepMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		query := models.TaskStepsQuery{
			TaskSteps: models.TaskSteps{
				UserID: helpers.ProcessUserIdFromContext(ctx),
				TaskID: ctx.Query("task_id"),
			},
		}
		trueVar := true
		falseVar := false

		if done := ctx.Query("done"); done != "" {
			switch done {
			case "true":
				query.IsDone = &trueVar
			case "false":
				query.IsDone = &falseVar
			}
		}
		if ctx.Query("done_at_lte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("done_at_lte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.DoneAtLte = &cal
		}
		if ctx.Query("done_at_gte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("done_at_gte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.DoneAtGte = &cal
		}
		if ctx.Query("time_due_lte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("time_due_lte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.TimeDueLte = &cal
		}
		if ctx.Query("time_due_gte") != "" {
			cal, err := helpers.ParseDate(ctx.Query("time_due_gte"))
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query.TimeDueGte = &cal
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
