package scopes

import (
	"habitask-backend-go/pkg/lib/structs"

	"gorm.io/gorm"
)

func ApplyPaginationAndSort(pagination *structs.PaginationAndSort) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination.SortBy != "" {
			if pagination.SortDirection == "" {
				pagination.SortDirection = "desc"
			}
			db = db.Order(pagination.SortBy + " " + pagination.SortDirection)
		}
		return db.Offset(pagination.GetOffset()).Limit(pagination.PerPage)
	}
}
