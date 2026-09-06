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

type TaskRepository interface {
	GetTasks(filter *models.TaskQuery, pagination *structs.PaginationAndSort) ([]*models.Task, error)
	GetTask(id string) (*models.Task, error)
	GetTaskSummary(filter *models.TaskQuery, pagination *structs.PaginationAndSort) ([]models.TaskSummary, error)
	CreateTask(task *models.Task) error
	UpdateTask(task *models.Task) error
	DeleteTask(id string) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) GetTasks(filter *models.TaskQuery, pagination *structs.PaginationAndSort) ([]*models.Task, error) {
	tasks := []*models.Task{}
	query := r.db.Model(&models.Task{}).Where(filter.Task)
	if filter.Search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.CompletedAtLte != nil {
		query = query.Where("completed_at <= ?", *filter.CompletedAtLte)
	}
	if filter.CompletedAtGte != nil {
		query = query.Where("completed_at >= ?", *filter.CompletedAtGte)
	}
	if filter.DueAtLte != nil {
		query = query.Where("due_at <= ?", *filter.DueAtLte)
	}
	if filter.DueAtGte != nil {
		query = query.Where("due_at >= ?", *filter.DueAtGte)
	}
	if filter.EstimatedDurationLte != nil {
		query = query.Where("estimated_duration <= ?", *filter.EstimatedDurationLte)
	}
	if filter.EstimatedDurationGte != nil {
		query = query.Where("estimated_duration >= ?", *filter.EstimatedDurationGte)
	}
	if filter.ActualDurationLte != nil {
		query = query.Where("actual_duration <= ?", *filter.ActualDurationLte)
	}
	if filter.ActualDurationGte != nil {
		query = query.Where("actual_duration >= ?", *filter.ActualDurationGte)
	}
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetTask(id string) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.First(task, id).Error
	return task, err
}

func (r *taskRepository) GetTaskSummary(filter *models.TaskQuery, pagination *structs.PaginationAndSort) ([]models.TaskSummary, error) {
	var tasks []models.TaskSummary
	groupExpr, err := helpers.GetSummaryExpression(helpers.SummaryGroup{
		GroupBy:         filter.GroupBy,
		DateGroup:       filter.DateGroup,
		Field:           filter.GroupBy,
		AllowedGrouping: taskSummaryGroups,
	})
	if err != nil {
		return nil, err
	}
	query := r.db.Model(&models.Task{}).Where(filter.Task).Select(
		fmt.Sprintf(`
		%s as label,
		COUNT(tasks.id) as count,
		SUM(tasks.is_completed) as completed_count,
		SUM(tasks.is_completed = 0) as incomplete_count,
		MAX(tasks.due_at) as due_count,
		MAX(tasks.completed_at) as last_task_date
		`, groupExpr),
	)
	if strings.HasPrefix(groupExpr, "users.") {
		query = query.Joins("JOIN users ON tasks.user_id = users.id")
	}
	if filter.Search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.CompletedAtLte != nil {
		query = query.Where("completed_at <= ?", *filter.CompletedAtLte)
	}
	if filter.CompletedAtGte != nil {
		query = query.Where("completed_at >= ?", *filter.CompletedAtGte)
	}
	if filter.DueAtLte != nil {
		query = query.Where("due_at <= ?", *filter.DueAtLte)
	}
	if filter.DueAtGte != nil {
		query = query.Where("due_at >= ?", *filter.DueAtGte)
	}
	if filter.EstimatedDurationLte != nil {
		query = query.Where("estimated_duration <= ?", *filter.EstimatedDurationLte)
	}
	if filter.EstimatedDurationGte != nil {
		query = query.Where("estimated_duration >= ?", *filter.EstimatedDurationGte)
	}
	if filter.ActualDurationLte != nil {
		query = query.Where("actual_duration <= ?", *filter.ActualDurationLte)
	}
	if filter.ActualDurationGte != nil {
		query = query.Where("actual_duration >= ?", *filter.ActualDurationGte)
	}
	err = query.Offset(pagination.GetOffset()).Limit(pagination.PerPage).Group(groupExpr).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *taskRepository) CreateTask(task *models.Task) error {
	err := r.db.Create(task).Error
	return err
}

func (r *taskRepository) UpdateTask(task *models.Task) error {
	err := r.db.Where("id = ?", task.ID).Updates(task).Error
	return err
}

func (r *taskRepository) DeleteTask(id string) error {
	task := &models.Task{}
	err := r.db.Where("id = ?", id).Delete(task).Error
	return err
}

var taskSummaryGroups = map[string]string{
	"created_at":   "tasks.created_at",
	"due_at":       "tasks.due_at",
	"name":         "tasks.name",
	"category":     "tasks.category",
	"is_completed": "tasks.is_completed",
	"completed_at": "tasks.completed_at",
	"user_id":      "tasks.user_id",
	"user_name":    "users.user_name",
}
