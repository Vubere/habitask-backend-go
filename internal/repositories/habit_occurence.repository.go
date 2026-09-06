package repositories

import (
	"fmt"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"
	"habitask-backend-go/pkg/lib/database/scopes"
	"habitask-backend-go/pkg/lib/structs"
	"strings"

	"gorm.io/gorm"
)

type HabitOccurenceRepository interface {
	GetHabitOccurences(filter *models.HabitOccurenceQuery, pagination *structs.PaginationAndSort) ([]*models.HabitOccurence, error)
	GetHabitOccurence(id string) (*models.HabitOccurence, error)
	GetHabitOccurenceSummary(filter *models.HabitOccurenceQuery, pagination *structs.PaginationAndSort) ([]models.HabitOccurenceSummary, error)
	CreateHabitOccurence(habitOccurence *models.HabitOccurence) error
	UpdateHabitOccurence(habitOccurence *models.HabitOccurence) error
	DeleteHabitOccurence(id string) error
}

type habitOccurenceRepository struct {
	db *gorm.DB
}

func NewHabitOccurenceRepository(db *gorm.DB) HabitOccurenceRepository {
	return &habitOccurenceRepository{db: db}
}

func (r *habitOccurenceRepository) GetHabitOccurences(filter *models.HabitOccurenceQuery, pagination *structs.PaginationAndSort) ([]*models.HabitOccurence, error) {
	habitOccurences := []*models.HabitOccurence{}
	query := r.db.Model(&models.HabitOccurence{}).Where(filter.HabitOccurence)
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.DateLte != nil {
		query = query.Where("date <= ?", *filter.DateLte)
	}
	if filter.DateGte != nil {
		query = query.Where("date >= ?", *filter.DateGte)
	}
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&habitOccurences).Error
	return habitOccurences, err
}

func (r *habitOccurenceRepository) GetHabitOccurence(id string) (*models.HabitOccurence, error) {
	habitOccurence := &models.HabitOccurence{}
	err := r.db.First(habitOccurence, id).Error
	return habitOccurence, err
}

func (r *habitOccurenceRepository) GetHabitOccurenceSummary(filter *models.HabitOccurenceQuery, pagination *structs.PaginationAndSort) ([]models.HabitOccurenceSummary, error) {
	var habitOccurences []models.HabitOccurenceSummary
	groupExpr, err := helpers.GetSummaryExpression(helpers.SummaryGroup{
		GroupBy:         filter.GroupBy,
		DateGroup:       filter.DateGroup,
		Field:           filter.GroupBy,
		AllowedGrouping: habitOccurenceSummaryGroups,
	})
	if err != nil {
		return nil, err
	}
	query := r.db.Model(&models.HabitOccurence{}).Where(filter.HabitOccurence).Select(
		fmt.Sprintf(`
		%s as label,
		COUNT(habit_occurences.id) as count,
		MAX(habit_occurences.date) as last_habit_date,
		habit_occurences.user_id
		`, groupExpr),
	)
	if strings.HasPrefix(groupExpr, "habits.") {
		query = query.Joins("JOIN habits ON habit_occurences.habit_id = habits.id")
	}
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.DateLte != nil {
		query = query.Where("date <= ?", *filter.DateLte)
	}
	if filter.DateGte != nil {
		query = query.Where("date >= ?", *filter.DateGte)
	}
	err = query.Offset(pagination.GetOffset()).Limit(pagination.PerPage).Group(groupExpr).Find(&habitOccurences).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *habitOccurenceRepository) CreateHabitOccurence(habitOccurence *models.HabitOccurence) error {
	err := r.db.Create(habitOccurence).Error
	return err
}

func (r *habitOccurenceRepository) UpdateHabitOccurence(habitOccurence *models.HabitOccurence) error {
	err := r.db.Where("id = ?", habitOccurence.ID).Updates(habitOccurence).Error
	return err
}

func (r *habitOccurenceRepository) DeleteHabitOccurence(id string) error {
	habitOccurence := &models.HabitOccurence{}
	err := r.db.Where("id = ?", id).Delete(habitOccurence).Error
	return err
}

var habitOccurenceSummaryGroups = map[string]string{
	"created_at":  "habit_occurences.created_at",
	"date":        "habit_occurences.date",
	"habit_id":    "habit_occurences.habit_id",
	"user_id":     "habit_occurences.user_id",
	"category":    "habits.category",
	"is_positive": "habits.is_positive",
	"habit_name":  "habits.habit_name",
}
