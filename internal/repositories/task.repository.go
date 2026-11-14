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

type TaskRepository interface {
	GetTasks(filter *models.TaskQuery, pagination *structs.PaginationAndSort) ([]models.Task, error)
	GetTask(id string) (*models.Task, error)
	GetTaskSummary(filter *models.TaskQuery, pagination *structs.PaginationAndSort) ([]models.TaskSummary, error)
	CreateTask(task *models.Task) error
	UpdateTask(id string, task *models.Task) error
	SetTaskArchived(id string, archived bool) error
	DeleteTask(id string) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) GetTasks(filter *models.TaskQuery, pagination *structs.PaginationAndSort) ([]models.Task, error) {
	tasks := []models.Task{}
	query := r.db.Model(&models.Task{}).Where(filter.Task)
	if filter.Archived == nil || !*filter.Archived {
		query = query.Where("tasks.archived_at IS NULL")
	}
	if filter.Search != "" {
		query = query.Where("tasks.title LIKE ? OR tasks.notes LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.DateLte != nil {
		query = query.Where("tasks.created_at <= ?", *filter.DateLte)
	}
	if filter.DateGte != nil {
		query = query.Where("tasks.created_at >= ?", *filter.DateGte)
	}
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetTask(id string) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.First(task, "id = ?", id).Error
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
		SUM(tasks.status = 'done') as completed_count,
		SUM(tasks.status = 'todo') as incomplete_count,
		SUM(tasks.due_date IS NOT NULL) as due_count,
		MAX(tasks.completed_at) as last_task_date
		`, groupExpr),
	)
	if strings.HasPrefix(groupExpr, "users.") {
		query = query.Joins("JOIN users ON tasks.user_id = users.id")
	}
	if filter.Search != "" {
		query = query.Where("tasks.title LIKE ? OR tasks.notes LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	err = query.Offset(pagination.GetOffset()).Limit(pagination.PerPage).Group(groupExpr).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) CreateTask(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) UpdateTask(id string, task *models.Task) error {
	return r.db.Model(&models.Task{}).Where("id = ?", id).Updates(task).Error
}

func (r *taskRepository) SetTaskArchived(id string, archived bool) error {
	if archived {
		return r.db.Model(&models.Task{}).Where("id = ?", id).Update("archived_at", time.Now()).Error
	}
	return r.db.Model(&models.Task{}).Where("id = ?", id).Update("archived_at", nil).Error
}

func (r *taskRepository) DeleteTask(id string) error {
	return r.SetTaskArchived(id, true)
}

var taskSummaryGroups = map[string]string{
	"created_at": "tasks.created_at",
	"due_date":   "tasks.due_date",
	"title":      "tasks.title",
	"priority":   "tasks.priority",
	"status":     "tasks.status",
	"completed_at": "tasks.completed_at",
	"user_id":      "tasks.user_id",
	"username":     "users.username",
}
