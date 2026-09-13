package middleware

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"

	"github.com/gin-gonic/gin"
)

func UserMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		query := models.UserQuery{}

		trueVar := true
		falseVar := false
		if ctx.Query("search") != "" {
			query.Search = ctx.Query("search")
		}
		if ctx.Query("active") != "" {
			query.Active = &trueVar
		}
		if ctx.Query("admin") != "" {
			adminFilter := ctx.Query("admin")
			switch adminFilter {
			case "true":
				query.Admin = &trueVar
			case "false":
				query.Admin = &falseVar
			}
		}
		if ctx.Query("username") != "" {
			query.Username = ctx.Query("username")
		}
		if ctx.Query("email") != "" {
			query.Email = ctx.Query("email")
		}
		if ctx.Query("first_name") != "" {
			query.FirstName = ctx.Query("first_name")
		}
		if ctx.Query("last_name") != "" {
			query.LastName = ctx.Query("last_name")
		}
		paginationAndSort := helpers.GetPaginationAndSortFromContext(ctx)
		ctx.Set("paginationAndSort", &paginationAndSort)
		ctx.Set("query", &query)
		ctx.Next()
	}
}
