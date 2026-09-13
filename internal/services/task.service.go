package services

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/pkg/lib/structs"
)

type TaskService interface {
	GetTasks(query models.TaskQuery, pagination *structs.PaginationAndSort) ([]models.Task, error)
	GetTaskSummary(query models.TaskQuery, pagination *structs.PaginationAndSort) ([]models.TaskSummary, error)
	GetTaskById(id string) (*models.Task, error)
	CreateTask(task *models.Task) error
	UpdateTask(id string, task *models.Task) error
	DeleteTask(id string) error
	ValidateEntry(task *models.Task) error
}

type taskService struct {
	repo     repositories.TaskRepository
	userRepo repositories.UserRepository
}

func NewTaskService(repo repositories.TaskRepository, userRepo repositories.UserRepository) TaskService {
	return &taskService{repo: repo, userRepo: userRepo}
}

func (s *taskService) GetTasks(query models.TaskQuery, pagination *structs.PaginationAndSort) ([]models.Task, error) {
	return s.repo.GetTasks(&query, pagination)
}

func (s *taskService) GetTaskSummary(query models.TaskQuery, pagination *structs.PaginationAndSort) ([]models.TaskSummary, error) {
	return s.repo.GetTaskSummary(&query, pagination)
}

func (s *taskService) GetTaskById(id string) (*models.Task, error) {
	return s.repo.GetTask(id)
}

func (s *taskService) CreateTask(task *models.Task) error {
	if err := s.ValidateEntry(task); err != nil {
		return err
	}
	return s.repo.CreateTask(task)
}

func (s *taskService) UpdateTask(id string, task *models.Task) error {
	if err := s.ValidateEntry(task); err != nil {
		return err
	}
	return s.repo.UpdateTask(id, task)
}

func (s *taskService) DeleteTask(id string) error {
	return s.repo.DeleteTask(id)
}

func (s *taskService) ValidateEntry(task *models.Task) error {
	if s.isForCreate(task) || task.UserID != "" {
		_, err := s.userRepo.GetUser(task.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *taskService) isForCreate(task *models.Task) bool {
	return task.ID == ""
}
