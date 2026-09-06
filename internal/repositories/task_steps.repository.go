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

type TaskStepsRepository interface {
	GetTaskSteps(filter *models.TaskStepsQuery, pagination *structs.PaginationAndSort) ([]*models.TaskSteps, error)
	GetTaskStep(id string) (*models.TaskSteps, error)
	GetTaskStepSummary(filter *models.TaskStepsQuery, pagination *structs.PaginationAndSort) ([]models.TaskStepSummary, error)
	CreateTaskStep(taskStep *models.TaskSteps) error
	UpdateTaskStep(taskStep *models.TaskSteps) error
	DeleteTaskStep(id string) error
}

type taskStepsRepository struct {
	db *gorm.DB
}

func NewTaskStepsRepository(db *gorm.DB) TaskStepsRepository {
	return &taskStepsRepository{db: db}
}

func (r *taskStepsRepository) GetTaskSteps(filter *models.TaskStepsQuery, pagination *structs.PaginationAndSort) ([]*models.TaskSteps, error) {
	taskSteps := []*models.TaskSteps{}
	query := r.db.Model(&models.TaskSteps{}).Where(filter.TaskSteps)
	if filter.Search != "" {
		query = query.Where("description LIKE ? OR time_due LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.DoneAtLte != nil {
		query = query.Where("done_at <= ?", *filter.DoneAtLte)
	}
	if filter.DoneAtGte != nil {
		query = query.Where("done_at >= ?", *filter.DoneAtGte)
	}
	if filter.TimeDueLte != nil {
		query = query.Where("time_due <= ?", *filter.TimeDueLte)
	}
	if filter.TimeDueGte != nil {
		query = query.Where("time_due >= ?", *filter.TimeDueGte)
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
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&taskSteps).Error
	return taskSteps, err
}

func (r *taskStepsRepository) GetTaskStep(id string) (*models.TaskSteps, error) {
	taskStep := &models.TaskSteps{}
	err := r.db.First(taskStep, id).Error
	return taskStep, err
}

func (r *taskStepsRepository) GetTaskStepSummary(filter *models.TaskStepsQuery, pagination *structs.PaginationAndSort) ([]models.TaskStepSummary, error) {
	var taskSteps []models.TaskStepSummary
	groupExpr, err := helpers.GetSummaryExpression(helpers.SummaryGroup{
		GroupBy:         filter.GroupBy,
		DateGroup:       filter.DateGroup,
		Field:           filter.GroupBy,
		AllowedGrouping: taskStepSummaryGroups,
	})
	if err != nil {
		return nil, err
	}
	query := r.db.Model(&models.TaskSteps{}).Where(filter.TaskSteps).Select(
		fmt.Sprintf(`
		%s as label,
		COUNT(task_steps.id) as count,
		SUM(task_steps.is_done) as done_count,
		MAX(task_steps.done_at) as last_task_step_date,
		MAX(task_steps.time_due) as time_due,
		SUM(task_steps.estimated_duration) as estimated_duration,
		SUM(task_steps.actual_duration) as actual_duration,
		SUM(task_steps.actual_duration) / SUM(task_steps.estimated_duration) as time_unit
		`, groupExpr),
	)
	if strings.HasPrefix(groupExpr, "users.") {
		query = query.Joins("JOIN users ON task_steps.user_id = users.id")
	}
	if filter.Search != "" {
		query = query.Where("description LIKE ? OR time_due LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.DoneAtLte != nil {
		query = query.Where("done_at <= ?", *filter.DoneAtLte)
	}
	if filter.DoneAtGte != nil {
		query = query.Where("done_at >= ?", *filter.DoneAtGte)
	}
	if filter.TimeDueLte != nil {
		query = query.Where("time_due <= ?", *filter.TimeDueLte)
	}
	if filter.TimeDueGte != nil {
		query = query.Where("time_due >= ?", *filter.TimeDueGte)
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
	err = query.Offset(pagination.GetOffset()).Limit(pagination.PerPage).Group(groupExpr).Find(&taskSteps).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *taskStepsRepository) CreateTaskStep(taskStep *models.TaskSteps) error {
	err := r.db.Create(taskStep).Error
	return err
}

func (r *taskStepsRepository) UpdateTaskStep(taskStep *models.TaskSteps) error {
	err := r.db.Where("id = ?", taskStep.ID).Updates(taskStep).Error
	return err
}

func (r *taskStepsRepository) DeleteTaskStep(id string) error {
	taskStep := &models.TaskSteps{}
	err := r.db.Where("id = ?", id).Delete(taskStep).Error
	return err
}

var taskStepSummaryGroups = map[string]string{
	"created_at":    "task_steps.created_at",
	"done_at":       "task_steps.done_at",
	"time_due":      "task_steps.time_due",
	"description":   "task_steps.description",
	"is_done":       "task_steps.is_done",
	"user_id":       "task_steps.user_id",
	"user_name":     "users.user_name",
	"task_id":       "task_steps.task_id",
	"task_name":     "tasks.name",
	"task_category": "tasks.category",
}
