package services

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/pkg/lib/structs"
)

type TaskStepService interface {
	GetTaskSteps(query models.TaskStepsQuery, pagination *structs.PaginationAndSort) ([]models.TaskSteps, error)
	GetTaskStepSummary(query models.TaskStepsQuery, pagination *structs.PaginationAndSort) ([]models.TaskStepSummary, error)
	GetTaskStepById(id string) (*models.TaskSteps, error)
	CreateTaskStep(taskStep *models.TaskSteps) error
	UpdateTaskStep(id string, taskStep *models.TaskSteps) error
	DeleteTaskStep(id string) error
	ValidateEntry(taskStep *models.TaskSteps) error
}

type taskStepService struct {
	repo     repositories.TaskStepsRepository
	userRepo repositories.UserRepository
	taskRepo repositories.TaskRepository
}

func NewTaskStepService(repo repositories.TaskStepsRepository, userRepo repositories.UserRepository, taskRepo repositories.TaskRepository) TaskStepService {
	return &taskStepService{repo: repo, userRepo: userRepo, taskRepo: taskRepo}
}

func (s *taskStepService) GetTaskSteps(query models.TaskStepsQuery, pagination *structs.PaginationAndSort) ([]models.TaskSteps, error) {
	return s.repo.GetTaskSteps(&query, pagination)
}

func (s *taskStepService) GetTaskStepSummary(query models.TaskStepsQuery, pagination *structs.PaginationAndSort) ([]models.TaskStepSummary, error) {
	return s.repo.GetTaskStepSummary(&query, pagination)
}

func (s *taskStepService) GetTaskStepById(id string) (*models.TaskSteps, error) {
	return s.repo.GetTaskStep(id)
}

func (s *taskStepService) CreateTaskStep(taskStep *models.TaskSteps) error {
	if err := s.ValidateEntry(taskStep); err != nil {
		return err
	}
	return s.repo.CreateTaskStep(taskStep)
}

func (s *taskStepService) UpdateTaskStep(id string, taskStep *models.TaskSteps) error {
	if err := s.ValidateEntry(taskStep); err != nil {
		return err
	}
	return s.repo.UpdateTaskStep(id, taskStep)
}

func (s *taskStepService) DeleteTaskStep(id string) error {
	return s.repo.DeleteTaskStep(id)
}

func (s *taskStepService) ValidateEntry(taskStep *models.TaskSteps) error {
	if s.isForCreate(taskStep) || taskStep.TaskID != "" {
		_, err := s.taskRepo.GetTask(taskStep.TaskID)
		if err != nil {
			return err
		}
	}
	if s.isForCreate(taskStep) || taskStep.UserID != "" {
		_, err := s.userRepo.GetUser(taskStep.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *taskStepService) isForCreate(taskStep *models.TaskSteps) bool {
	return taskStep.ID == ""
}
