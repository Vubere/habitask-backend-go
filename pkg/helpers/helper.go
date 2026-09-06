package helpers

import "fmt"

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
