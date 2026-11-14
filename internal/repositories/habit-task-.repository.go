package repositories

import (
	"encoding/json"
	"sort"
	"time"

	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/lib/structs"

	"gorm.io/gorm"
)

type HabitTaskRepository interface {
	GetHabitTaskDisplayDetails(filter dtos.HabitTaskFilterDTO, pagination *structs.PaginationAndSort, search string) ([]dtos.HabitTaskDTO, error)
}

type habitTaskRepository struct {
	db *gorm.DB
}

func NewHabitTaskRepository(db *gorm.DB) HabitTaskRepository {
	return &habitTaskRepository{db: db}
}

func (r *habitTaskRepository) GetHabitTaskDisplayDetails(
	filter dtos.HabitTaskFilterDTO,
	pagination *structs.PaginationAndSort,
	search string,
) ([]dtos.HabitTaskDTO, error) {
	rows := []dtos.HabitTaskDTO{}
	const layout = "2006-01-02 15:04:05"

	if filter.Type == "" || filter.Type == "habit" {
		habits := []models.Habit{}
		query := r.db.Model(&models.Habit{}).
			Where("habits.user_id = ?", filter.UserID).
			Where("habits.archived_at IS NULL")
		if filter.CreatedAtGTE != nil && !filter.CreatedAtGTE.IsZero() {
			query = query.Where("habits.created_at >= ?", *filter.CreatedAtGTE)
		}
		if filter.CreatedAtLTE != nil && !filter.CreatedAtLTE.IsZero() {
			query = query.Where("habits.created_at <= ?", *filter.CreatedAtLTE)
		}
		if search != "" {
			like := "%" + search + "%"
			query = query.Where("habits.title LIKE ? OR habits.notes LIKE ?", like, like)
		}
		if err := query.Find(&habits).Error; err != nil {
			return nil, err
		}

		progress := r.todayProgressByHabit(habits)
		streaks := r.streaksByHabit(habits)
		for i := range habits {
			h := habits[i]
			streak := streaks[h.ID]
			rows = append(rows, dtos.HabitTaskDTO{
				Id:        h.ID,
				Type:      "habit",
				Title:     h.Title,
				Notes:     h.Notes,
				Icon:      h.Icon,
				Color:     h.Color,
				Status:    "",
				Progress:  progress[h.ID],
				Target:    h.TargetPerDay,
				Streak:    &streak,
				CreatedAt: h.CreatedAt.Format(layout),
				UpdatedAt: h.UpdatedAt.Format(layout),
			})
		}
	}

	if filter.Type == "" || filter.Type == "task" {
		tasks := []models.Task{}
		query := r.db.Model(&models.Task{}).
			Where("tasks.user_id = ?", filter.UserID).
			Where("tasks.archived_at IS NULL")
		if filter.CreatedAtGTE != nil && !filter.CreatedAtGTE.IsZero() {
			query = query.Where("tasks.created_at >= ?", *filter.CreatedAtGTE)
		}
		if filter.CreatedAtLTE != nil && !filter.CreatedAtLTE.IsZero() {
			query = query.Where("tasks.created_at <= ?", *filter.CreatedAtLTE)
		}
		if search != "" {
			like := "%" + search + "%"
			query = query.Where("tasks.title LIKE ? OR tasks.notes LIKE ?", like, like)
		}
		if err := query.Find(&tasks).Error; err != nil {
			return nil, err
		}

		for i := range tasks {
			t := tasks[i]
			progress, target := subtaskProgress(t.Subtasks)
			rows = append(rows, dtos.HabitTaskDTO{
				Id:        t.ID,
				Type:      "task",
				Title:     t.Title,
				Notes:     t.Notes,
				Icon:      t.Icon,
				Color:     t.Color,
				Status:    string(t.Status),
				Progress:  progress,
				Target:    target,
				CreatedAt: t.CreatedAt.Format(layout),
				UpdatedAt: t.UpdatedAt.Format(layout),
			})
		}
	}

	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].CreatedAt > rows[j].CreatedAt
	})

	start := pagination.GetOffset()
	if start > len(rows) {
		start = len(rows)
	}
	end := start + pagination.PerPage
	if end > len(rows) {
		end = len(rows)
	}
	return rows[start:end], nil
}

func (r *habitTaskRepository) todayProgressByHabit(habits []models.Habit) map[string]int {
	progress := map[string]int{}
	if len(habits) == 0 {
		return progress
	}
	ids := make([]string, 0, len(habits))
	for i := range habits {
		ids = append(ids, habits[i].ID)
	}
	type aggRow struct {
		HabitID string
		Total   int
	}
	var agg []aggRow
	r.db.Model(&models.HabitCompletion{}).
		Select("habit_id, SUM(count) as total").
		Where("habit_id IN ?", ids).
		Where("date = ?", time.Now().Format("2006-01-02")).
		Group("habit_id").
		Scan(&agg)
	for _, a := range agg {
		progress[a.HabitID] = a.Total
	}
	return progress
}

func (r *habitTaskRepository) streaksByHabit(habits []models.Habit) map[string]int {
	streaks := map[string]int{}
	if len(habits) == 0 {
		return streaks
	}
	ids := make([]string, 0, len(habits))
	for i := range habits {
		ids = append(ids, habits[i].ID)
	}
	completions := []models.HabitCompletion{}
	if err := r.db.Where("habit_id IN ?", ids).Order("date asc").Find(&completions).Error; err != nil {
		return streaks
	}
	byHabit := map[string][]models.HabitCompletion{}
	for _, c := range completions {
		byHabit[c.HabitID] = append(byHabit[c.HabitID], c)
	}
	for i := range habits {
		h := habits[i]
		stats := computeHabitStats(&h, byHabit[h.ID])
		streaks[h.ID] = stats.CurrentStreak
	}
	return streaks
}

func subtaskProgress(raw string) (int, int) {
	if raw == "" {
		return 0, 0
	}
	var subtasks []struct {
		Done bool `json:"done"`
	}
	if err := json.Unmarshal([]byte(raw), &subtasks); err != nil {
		return 0, 0
	}
	done := 0
	for _, s := range subtasks {
		if s.Done {
			done++
		}
	}
	return done, len(subtasks)
}
