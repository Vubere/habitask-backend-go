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

type HabitRepository interface {
	GetHabits(filter *models.HabitQuery, pagination *structs.PaginationAndSort) ([]*models.Habit, error)
	GetHabit(id string) (*models.Habit, error)
	GetHabitSummary(filter *models.HabitQuery, pagination *structs.PaginationAndSort) ([]models.HabbitSummary, error)
	CreateHabit(habit *models.Habit) error
	UpdateHabit(habit *models.Habit) error
	DeleteHabit(id string) error
}

type habitRepository struct {
	db *gorm.DB
}

func NewHabitRepository(db *gorm.DB) HabitRepository {
	return &habitRepository{db: db}
}

func (r *habitRepository) GetHabits(filter *models.HabitQuery, pagination *structs.PaginationAndSort) ([]*models.Habit, error) {
	habits := []*models.Habit{}
	query := r.db.Model(&models.Habit{}).Where(filter.Habit)
	if filter.Search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.LastDoneLte != nil {
		query = query.Where("last_done <= ?", *filter.LastDoneLte)
	}
	if filter.LastDoneGte != nil {
		query = query.Where("last_done >= ?", *filter.LastDoneGte)
	}
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&habits).Error
	return habits, err
}

func (r *habitRepository) GetHabit(id string) (*models.Habit, error) {
	habit := &models.Habit{}
	err := r.db.First(habit, id).Error
	return habit, err
}

func (r *habitRepository) GetHabitSummary(filter *models.HabitQuery, pagination *structs.PaginationAndSort) ([]models.HabbitSummary, error) {
	var habits []models.HabbitSummary
	groupExpr, err := helpers.GetSummaryExpression(helpers.SummaryGroup{
		GroupBy:         filter.GroupBy,
		DateGroup:       filter.DateGroup,
		Field:           filter.GroupBy,
		AllowedGrouping: habitSummaryGroups,
	})
	if err != nil {
		return nil, err
	}
	query := r.db.Model(&models.Habit{}).Where(filter.Habit).Select(
		fmt.Sprintf(`
		%s as label,
		COUNT(habits.id) as count,
		SUM(habits.is_positive) as positive_count,
		SUM(habits.is_positive = 0) as negative_count,
		MAX(habits.last_done) as last_habbit_date
		`, groupExpr),
	)
	if strings.HasPrefix(groupExpr, "users.") {
		query = query.Joins("JOIN users ON habits.user_id = users.id")
	}
	if filter.Search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.LastDoneLte != nil {
		query = query.Where("last_done <= ?", *filter.LastDoneLte)
	}
	if filter.LastDoneGte != nil {
		query = query.Where("last_done >= ?", *filter.LastDoneGte)
	}
	err = query.Offset(pagination.GetOffset()).Limit(pagination.PerPage).Group(groupExpr).Find(&habits).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *habitRepository) CreateHabit(habit *models.Habit) error {
	err := r.db.Create(habit).Error
	return err
}

func (r *habitRepository) UpdateHabit(habit *models.Habit) error {
	err := r.db.Where("id = ?", habit.ID).Updates(habit).Error
	return err
}

func (r *habitRepository) DeleteHabit(id string) error {
	habit := &models.Habit{}
	err := r.db.Where("id = ?", id).Delete(habit).Error
	return err
}

var habitSummaryGroups = map[string]string{
	"created_at": "habits.created_at",
	"last_done":  "habits.last_done",
	"name":       "habits.name",
	"category":   "habits.category",
	"positive":   "habits.is_positive",
	"user_id":    "habits.user_id",
	"username":   "users.username",
}
