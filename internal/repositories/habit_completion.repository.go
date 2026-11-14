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

type HabitCompletionRepository interface {
	GetHabitCompletions(filter *models.HabitCompletionQuery, pagination *structs.PaginationAndSort) ([]models.HabitCompletion, error)
	GetHabitCompletion(id string) (*models.HabitCompletion, error)
	GetHabitCompletionSummary(filter *models.HabitCompletionQuery, pagination *structs.PaginationAndSort) ([]models.HabitCompletionSummary, error)
	CreateHabitCompletion(completion *models.HabitCompletion) error
	UpsertHabitCompletion(completion *models.HabitCompletion) error
	UpdateHabitCompletion(id string, completion *models.HabitCompletion) error
	DeleteHabitCompletion(id string) error
}

type habitCompletionRepository struct {
	db *gorm.DB
}

func NewHabitCompletionRepository(db *gorm.DB) HabitCompletionRepository {
	return &habitCompletionRepository{db: db}
}

func (r *habitCompletionRepository) GetHabitCompletions(filter *models.HabitCompletionQuery, pagination *structs.PaginationAndSort) ([]models.HabitCompletion, error) {
	completions := []models.HabitCompletion{}
	query := r.db.Model(&models.HabitCompletion{}).Where(filter.HabitCompletion)
	if filter.DateLte != nil {
		query = query.Where("date <= ?", *filter.DateLte)
	}
	if filter.DateGte != nil {
		query = query.Where("date >= ?", *filter.DateGte)
	}
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&completions).Error
	return completions, err
}

func (r *habitCompletionRepository) GetHabitCompletion(id string) (*models.HabitCompletion, error) {
	completion := &models.HabitCompletion{}
	err := r.db.First(completion, "id = ?", id).Error
	return completion, err
}

func (r *habitCompletionRepository) GetHabitCompletionSummary(filter *models.HabitCompletionQuery, pagination *structs.PaginationAndSort) ([]models.HabitCompletionSummary, error) {
	var summaries []models.HabitCompletionSummary
	groupExpr, err := helpers.GetSummaryExpression(helpers.SummaryGroup{
		GroupBy:         filter.GroupBy,
		DateGroup:       filter.DateGroup,
		Field:           filter.GroupBy,
		AllowedGrouping: habitCompletionSummaryGroups,
	})
	if err != nil {
		return nil, err
	}
	query := r.db.Model(&models.HabitCompletion{}).Where(filter.HabitCompletion).Select(
		fmt.Sprintf(`
		%s as label,
		COUNT(habit_completions.id) as count,
		SUM(habit_completions.count) as total_count,
		MAX(habit_completions.date) as last_habit_date
		`, groupExpr),
	)
	if strings.HasPrefix(groupExpr, "habits.") {
		query = query.Joins("JOIN habits ON habit_completions.habit_id = habits.id")
	}
	err = query.Offset(pagination.GetOffset()).Limit(pagination.PerPage).Group(groupExpr).Find(&summaries).Error
	return summaries, err
}

func (r *habitCompletionRepository) CreateHabitCompletion(completion *models.HabitCompletion) error {
	return r.db.Create(completion).Error
}

func (r *habitCompletionRepository) UpsertHabitCompletion(completion *models.HabitCompletion) error {
	existing := &models.HabitCompletion{}
	date := completion.Date.Format("2006-01-02")
	err := r.db.Where("habit_id = ? AND date = ?", completion.HabitID, date).First(existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(completion).Error
	}
	if err != nil {
		return err
	}
	existing.Count = completion.Count
	existing.CompletedAt = completion.CompletedAt
	return r.db.Save(existing).Error
}

func (r *habitCompletionRepository) UpdateHabitCompletion(id string, completion *models.HabitCompletion) error {
	return r.db.Model(&models.HabitCompletion{}).Where("id = ?", id).Updates(completion).Error
}

func (r *habitCompletionRepository) DeleteHabitCompletion(id string) error {
	completion := &models.HabitCompletion{}
	return r.db.Where("id = ?", id).Delete(completion).Error
}

var habitCompletionSummaryGroups = map[string]string{
	"created_at":  "habit_completions.created_at",
	"date":        "habit_completions.date",
	"habit_id":    "habit_completions.habit_id",
	"user_id":     "habit_completions.user_id",
	"habit_title": "habits.title",
}
