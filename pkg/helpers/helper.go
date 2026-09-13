package helpers

import (
	"errors"
	"fmt"
	"habitask-backend-go/pkg/lib/structs"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SummaryGroup struct {
	GroupBy         string
	DateGroup       string
	Field           string
	AllowedGrouping map[string]string
}

func GetSummaryExpression(summaryGroup SummaryGroup) (string, error) {
	group, ok := summaryGroup.AllowedGrouping[summaryGroup.GroupBy]
	if !ok {
		return "", fmt.Errorf("invalid group by: %s", summaryGroup.GroupBy)
	}
	if summaryGroup.DateGroup != "" {
		return GetSummaryDateExpression(summaryGroup.DateGroup, summaryGroup.GroupBy)
	}
	return group, nil
}

func GetSummaryDateExpression(dateGroup string, field string) (string, error) {
	if field == "" {
		return "", fmt.Errorf("date group field cannot be empty")
	}
	switch dateGroup {

	case "day":
		return "DATE_FORMAT(" + field + ", '%Y-%m-%d')", nil

	case "week":
		return "DATE_FORMAT(" + field + ", '%x-W%v')", nil

	case "week_day":
		return "DATE_SUB(DATE(" + field + "), INTERVAL DAYOFWEEK(" + field + ") DAY)", nil

	case "month":
		return "DATE_FORMAT(" + field + ", '%Y-%m')", nil

	case "month_name":
		return "DATE_FORMAT(" + field + ", '%Y-%M')", nil

	case "year":
		return "DATE_FORMAT(" + field + ", '%Y')", nil

	default:
		return "", fmt.Errorf("invalid date group: %s", dateGroup)
	}
}

func GetPaginationAndSortFromContext(ctx *gin.Context) (paginationAndSort structs.PaginationAndSort) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(ctx.Query("per_page"))
	if err != nil {
		perPage = 10
	}
	paginationAndSort.Page = page
	paginationAndSort.PerPage = perPage
	sortBy := ctx.Query("sort_by")
	sortDirection := ctx.Query("sort_direction")
	if sortBy != "" && sortDirection != "" {
		paginationAndSort.SortBy = sortBy
		paginationAndSort.SortDirection = sortDirection
	}
	if sortBy == "" && sortDirection == "" {
		paginationAndSort.SortBy = "created_at"
		paginationAndSort.SortDirection = "desc"
	}
	return paginationAndSort
}

func ProcessUserIdFromContext(ctx *gin.Context) string {
	userId := ctx.Query("user_id")
	if ctx.GetString("role") != "admin" {
		userId = ctx.GetString("userId")
	}
	return userId
}

func ParseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("no date strning passed")
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
	}

	var err error

	for _, layout := range layouts {
		var parsed time.Time
		parsed, err = time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, err
}
