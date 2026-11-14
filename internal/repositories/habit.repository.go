package repositories

import (
	"fmt"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"
	"habitask-backend-go/pkg/lib/database/scopes"
	"habitask-backend-go/pkg/lib/structs"
	"strings"
	"time"

	"gorm.io/gorm"
)

type HabitRepository interface {
	GetHabits(filter *models.HabitQuery, pagination *structs.PaginationAndSort) ([]models.Habit, error)
	GetHabit(id string) (*models.Habit, error)
	GetHabitSummary(filter *models.HabitQuery, pagination *structs.PaginationAndSort) ([]models.HabitSummary, error)
	GetHabitStats(id string) (*models.HabitStats, error)
	CreateHabit(habit *models.Habit) error
	UpdateHabit(id string, habit *models.Habit) error
	SetHabitArchived(id string, archived bool) error
	DeleteHabit(id string) error
}

type habitRepository struct {
	db *gorm.DB
}

func NewHabitRepository(db *gorm.DB) HabitRepository {
	return &habitRepository{db: db}
}

func (r *habitRepository) GetHabits(filter *models.HabitQuery, pagination *structs.PaginationAndSort) ([]models.Habit, error) {
	habits := []models.Habit{}
	query := r.db.Model(&models.Habit{}).Where(filter.Habit)
	if filter.Archived == nil || !*filter.Archived {
		query = query.Where("habits.archived_at IS NULL")
	}
	if filter.Search != "" {
		query = query.Where("habits.title LIKE ? OR habits.notes LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.DateLte != nil {
		query = query.Where("habits.created_at <= ?", *filter.DateLte)
	}
	if filter.DateGte != nil {
		query = query.Where("habits.created_at >= ?", *filter.DateGte)
	}
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&habits).Error
	return habits, err
}

func (r *habitRepository) GetHabit(id string) (*models.Habit, error) {
	habit := &models.Habit{}
	err := r.db.First(habit, "id = ?", id).Error
	return habit, err
}

func (r *habitRepository) GetHabitSummary(filter *models.HabitQuery, pagination *structs.PaginationAndSort) ([]models.HabitSummary, error) {
	var habits []models.HabitSummary
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
		SUM(habits.polarity = 'positive') as positive_count,
		SUM(habits.polarity = 'negative') as negative_count,
		SUM(habits.target_per_day) as target_total
		`, groupExpr),
	)
	if strings.HasPrefix(groupExpr, "users.") {
		query = query.Joins("JOIN users ON habits.user_id = users.id")
	}
	if filter.Search != "" {
		query = query.Where("habits.title LIKE ? OR habits.notes LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	err = query.Offset(pagination.GetOffset()).Limit(pagination.PerPage).Group(groupExpr).Find(&habits).Error
	return habits, err
}

func (r *habitRepository) GetHabitStats(id string) (*models.HabitStats, error) {
	habit := &models.Habit{}
	if err := r.db.First(habit, "id = ?", id).Error; err != nil {
		return nil, err
	}
	var completions []models.HabitCompletion
	if err := r.db.Where("habit_id = ?", id).Order("date asc").Find(&completions).Error; err != nil {
		return nil, err
	}
	return computeHabitStats(habit, completions), nil
}

func (r *habitRepository) CreateHabit(habit *models.Habit) error {
	return r.db.Create(habit).Error
}

func (r *habitRepository) UpdateHabit(id string, habit *models.Habit) error {
	return r.db.Model(&models.Habit{}).Where("id = ?", id).Updates(habit).Error
}

func (r *habitRepository) SetHabitArchived(id string, archived bool) error {
	if archived {
		return r.db.Model(&models.Habit{}).Where("id = ?", id).Update("archived_at", time.Now()).Error
	}
	return r.db.Model(&models.Habit{}).Where("id = ?", id).Update("archived_at", nil).Error
}

func (r *habitRepository) DeleteHabit(id string) error {
	return r.SetHabitArchived(id, true)
}

func computeHabitStats(habit *models.Habit, completions []models.HabitCompletion) *models.HabitStats {
	perDay := map[string]int{}
	for _, c := range completions {
		perDay[c.Date.Format("2006-01-02")] += c.Count
	}
	target := habit.TargetPerDay
	if target < 1 {
		target = 1
	}
	start := habit.StartDate
	if start.IsZero() {
		start = time.Now()
	}
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	today := time.Now().UTC()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	isComplete := func(day time.Time) bool {
		return perDay[day.Format("2006-01-02")] >= target
	}

	best, run, total, done := 0, 0, 0, 0
	for d := start; !d.After(today); d = d.AddDate(0, 0, 1) {
		total++
		if isComplete(d) {
			run++
			done++
			if run > best {
				best = run
			}
		} else {
			run = 0
		}
	}

	current := 0
	day := today
	if !isComplete(day) {
		day = day.AddDate(0, 0, -1)
	}
	for !day.Before(start) && isComplete(day) {
		current++
		day = day.AddDate(0, 0, -1)
	}

	rate := 0.0
	if total > 0 {
		rate = float64(done) / float64(total)
	}

	return &models.HabitStats{
		CurrentStreak:  current,
		BestStreak:     best,
		CompletionRate: rate,
	}
}

var habitSummaryGroups = map[string]string{
	"created_at": "habits.created_at",
	"start_date": "habits.start_date",
	"title":      "habits.title",
	"polarity":   "habits.polarity",
	"user_id":    "habits.user_id",
	"username":   "users.username",
}
