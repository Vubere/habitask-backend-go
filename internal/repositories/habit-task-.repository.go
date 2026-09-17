package repositories

import (
	"habitask-backend-go/internal/dtos"
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
	var habitTasks []dtos.HabitTaskDTO

	query := `
		SELECT *
		FROM (
			SELECT 
				habits.id AS id,
				habits.name AS name,
				habits.description AS description,
				habits.category AS category,
				habits.status AS status,
				habits.progress AS progress,
				habits.streak AS streak,
				habits.created_at AS created_at,
				habits.updated_at AS updated_at
			FROM habits
			WHERE habits.user_id = ?

			UNION ALL

			SELECT 
				tasks.id AS id,
				tasks.name AS name,
				tasks.description AS description,
				tasks.category AS category,
				tasks.status AS status,
				tasks.progress AS progress,
				NULL AS streak,
				tasks.created_at AS created_at,
				tasks.updated_at AS updated_at
			FROM tasks
			WHERE tasks.user_id = ?
		) AS habit_tasks
		OFFSET ?
		LIMIT ?
		WHERE 1 = 1
	`

	args := []interface{}{
		filter.UserID,
		filter.UserID,
		pagination.GetOffset(),
		pagination.PerPage,
	}

	if filter.CreatedAtGTE != nil && filter.CreatedAtGTE.IsZero() == false {
		query += ` AND habit_tasks.created_at >= ?`
		args = append(args, filter.CreatedAtGTE)
	}

	if filter.CreatedAtLTE != nil && filter.CreatedAtLTE.IsZero() == false {
		query += ` AND habit_tasks.created_at <= ?`
		args = append(args, filter.CreatedAtLTE)
	}

	if search != "" {
		query += `
			AND (
				habit_tasks.name LIKE ?
				OR habit_tasks.description LIKE ?
				OR habit_tasks.category LIKE ?
			)
		`

		searchValue := "%" + search + "%"

		args = append(
			args,
			searchValue,
			searchValue,
			searchValue,
		)
	}

	query += ` ORDER BY habit_tasks.created_at DESC`

	err := r.db.Raw(query, args...).Scan(&habitTasks).Error

	return habitTasks, err
}
