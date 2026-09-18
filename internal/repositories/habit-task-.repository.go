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
				(SELECT COUNT(*) FROM habit_actions WHERE habit_actions.habit_id = habits.id AND habit_actions.is_done = 1) AS progress,
				(SELECT COUNT(*) FROM habit_actions WHERE habit_actions.habit_id = habits.id) AS steps,
				habits.streak AS streak,
				'habit' AS type,
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
				(SELECT COUNT(*) FROM task_steps WHERE task_steps.task_id = tasks.id AND task_steps.is_done = 1) AS progress,
				(SELECT COUNT(*) FROM task_steps WHERE task_steps.task_id = tasks.id) AS steps,
				NULL AS streak,
				'task' AS type,
				tasks.created_at AS created_at,
				tasks.updated_at AS updated_at
			FROM tasks
			WHERE tasks.user_id = ?
		) AS habit_tasks
		WHERE 1 = 1
	
	`

	args := []interface{}{
		filter.UserID,
		filter.UserID,
	}

	if filter.CreatedAtGTE != nil && filter.CreatedAtGTE.IsZero() == false {
		query += ` AND habit_tasks.created_at >= ?`
		args = append(args, filter.CreatedAtGTE)
	}

	if filter.CreatedAtLTE != nil && filter.CreatedAtLTE.IsZero() == false {
		query += ` AND habit_tasks.created_at <= ?`
		args = append(args, filter.CreatedAtLTE)
	}
	if filter.Type != "" {
		query += ` AND habit_tasks.type = ?`
		args = append(args, filter.Type)
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

	query += `
		ORDER BY habit_tasks.created_at DESC
		LIMIT ?
		OFFSET ?
	`
	args = append(args, pagination.PerPage, pagination.GetOffset())

	err := r.db.Raw(query, args...).Scan(&habitTasks).Error

	return habitTasks, err
}
