package middleware

import (
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/pkg/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HabitTaskMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := dtos.HabitTaskFilterDTO{
			UserID: c.GetString("user_id"),
		}
		if c.Query("created_at_lte") != "" {
			cal, err := helpers.ParseDate(c.Query("created_at_lte"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			query.CreatedAtLTE = &cal
		}
		if c.Query("created_at_gte") != "" {
			cal, err := helpers.ParseDate(c.Query("created_at_gte"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			query.CreatedAtGTE = &cal
		}
		if itemType := c.Query("type"); itemType != "" {
			query.Type = itemType
		}
		pagination := helpers.GetPaginationAndSortFromContext(c)
		c.Set("pagination", &pagination)
		c.Set("query", &query)
		c.Next()
	}
}
