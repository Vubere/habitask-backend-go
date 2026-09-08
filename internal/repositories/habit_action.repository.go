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

type HabitActionRepository interface {
	GetHabitActions(filter *models.HabitActionQuery, pagination *structs.PaginationAndSort) ([]*models.HabitAction, error)
	GetHabitAction(id string) (*models.HabitAction, error)
	GetHabitActionSummary(filter *models.HabitActionQuery, pagination *structs.PaginationAndSort) ([]models.HabitActionSummary, error)
	CreateHabitAction(habitAction *models.HabitAction) error
	UpdateHabitAction(habitAction *models.HabitAction) error
	DeleteHabitAction(id string) error
}

type habitActionRepository struct {
	db *gorm.DB
}

func NewHabitActionRepository(db *gorm.DB) HabitActionRepository {
	return &habitActionRepository{db: db}
}

func (r *habitActionRepository) GetHabitActions(filter *models.HabitActionQuery, pagination *structs.PaginationAndSort) ([]*models.HabitAction, error) {
	habitActions := []*models.HabitAction{}
	query := r.db.Model(&models.HabitAction{}).Where(filter.HabitAction)
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.DateLte != nil {
		query = query.Where("date <= ?", *filter.DateLte)
	}
	if filter.DateGte != nil {
		query = query.Where("date >= ?", *filter.DateGte)
	}
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&habitActions).Error
	return habitActions, err
}

func (r *habitActionRepository) GetHabitAction(id string) (*models.HabitAction, error) {
	habitAction := &models.HabitAction{}
	err := r.db.First(habitAction, id).Error
	return habitAction, err
}

func (r *habitActionRepository) GetHabitActionSummary(filter *models.HabitActionQuery, pagination *structs.PaginationAndSort) ([]models.HabitActionSummary, error) {
	var habitActions []models.HabitActionSummary
	groupExpr, err := helpers.GetSummaryExpression(helpers.SummaryGroup{
		GroupBy:         filter.GroupBy,
		DateGroup:       filter.DateGroup,
		Field:           filter.GroupBy,
		AllowedGrouping: habitActionSummaryGroups,
	})
	if err != nil {
		return nil, err
	}
	query := r.db.Model(&models.HabitAction{}).Where(filter.HabitAction).Select(
		fmt.Sprintf(`
		%s as label,
		COUNT(habit_actions.id) as count,
		SUM(habit_actions.cost_incurred) as cost_incurred
		`, groupExpr),
	)
	if strings.HasPrefix(groupExpr, "habits.") {
		query = query.Joins("JOIN habits ON habit_actions.habit_id = habits.id")
	}
	if strings.HasPrefix(groupExpr, "users.") {
		query = query.Joins("JOIN users ON habit_actions.user_id = users.id")
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
	err = query.Offset(pagination.GetOffset()).Limit(pagination.PerPage).Group(groupExpr).Find(&habitActions).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *habitActionRepository) CreateHabitAction(habitAction *models.HabitAction) error {
	err := r.db.Create(habitAction).Error
	return err
}

func (r *habitActionRepository) UpdateHabitAction(habitAction *models.HabitAction) error {
	err := r.db.Where("id = ?", habitAction.ID).Updates(habitAction).Error
	return err
}

func (r *habitActionRepository) DeleteHabitAction(id string) error {
	habitAction := &models.HabitAction{}
	err := r.db.Where("id = ?", id).Delete(habitAction).Error
	return err
}

var habitActionSummaryGroups = map[string]string{
	"created_at":    "habit_actions.created_at",
	"date":          "habit_actions.date",
	"title":         "habit_actions.title",
	"category":      "habit_actions.category",
	"cost_incurred": "habit_actions.cost_incurred",
	"is_positive":   "habit_actions.is_positive",
	"user_id":       "habit_actions.user_id",
	"username":      "users.username",
	"habit_id":      "habit_actions.habit_id",
	"habit_name":    "habits.habit_name",
}
