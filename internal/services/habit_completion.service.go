package services

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/pkg/lib/structs"
)

type HabitCompletionService interface {
	GetHabitCompletions(query models.HabitCompletionQuery, pagination *structs.PaginationAndSort) ([]models.HabitCompletion, error)
	GetHabitCompletionSummary(query models.HabitCompletionQuery, pagination *structs.PaginationAndSort) ([]models.HabitCompletionSummary, error)
	GetHabitCompletionById(id string) (*models.HabitCompletion, error)
	CreateHabitCompletion(completion *models.HabitCompletion) error
	LogHabitCompletion(completion *models.HabitCompletion) error
	UpdateHabitCompletion(id string, completion *models.HabitCompletion) error
	DeleteHabitCompletion(id string) error
	ValidateEntry(completion *models.HabitCompletion) error
}

type habitCompletionService struct {
	repo      repositories.HabitCompletionRepository
	habitRepo repositories.HabitRepository
	userRepo  repositories.UserRepository
}

func NewHabitCompletionService(
	repo repositories.HabitCompletionRepository,
	habitRepo repositories.HabitRepository,
	userRepo repositories.UserRepository,
) HabitCompletionService {
	return &habitCompletionService{repo: repo, habitRepo: habitRepo, userRepo: userRepo}
}

func (s *habitCompletionService) GetHabitCompletions(query models.HabitCompletionQuery, pagination *structs.PaginationAndSort) ([]models.HabitCompletion, error) {
	return s.repo.GetHabitCompletions(&query, pagination)
}

func (s *habitCompletionService) GetHabitCompletionSummary(query models.HabitCompletionQuery, pagination *structs.PaginationAndSort) ([]models.HabitCompletionSummary, error) {
	return s.repo.GetHabitCompletionSummary(&query, pagination)
}

func (s *habitCompletionService) GetHabitCompletionById(id string) (*models.HabitCompletion, error) {
	return s.repo.GetHabitCompletion(id)
}

func (s *habitCompletionService) CreateHabitCompletion(completion *models.HabitCompletion) error {
	if err := s.ValidateEntry(completion); err != nil {
		return err
	}
	return s.repo.CreateHabitCompletion(completion)
}

func (s *habitCompletionService) LogHabitCompletion(completion *models.HabitCompletion) error {
	if err := s.ValidateEntry(completion); err != nil {
		return err
	}
	return s.repo.UpsertHabitCompletion(completion)
}

func (s *habitCompletionService) UpdateHabitCompletion(id string, completion *models.HabitCompletion) error {
	if err := s.ValidateEntry(completion); err != nil {
		return err
	}
	return s.repo.UpdateHabitCompletion(id, completion)
}

func (s *habitCompletionService) DeleteHabitCompletion(id string) error {
	return s.repo.DeleteHabitCompletion(id)
}

func (s *habitCompletionService) ValidateEntry(completion *models.HabitCompletion) error {
	if s.isForCreate(completion) || completion.HabitID != "" {
		_, err := s.habitRepo.GetHabit(completion.HabitID)
		if err != nil {
			return err
		}
	}
	if s.isForCreate(completion) || completion.UserID != "" {
		_, err := s.userRepo.GetUser(completion.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *habitCompletionService) isForCreate(completion *models.HabitCompletion) bool {
	return completion.ID == ""
}
