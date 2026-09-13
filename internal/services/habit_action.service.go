package services

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/pkg/lib/structs"
)

type HabitActionService interface {
	GetHabitActions(query models.HabitActionQuery, pagination *structs.PaginationAndSort) ([]models.HabitAction, error)
	GetHabitActionSummary(query models.HabitActionQuery, pagination *structs.PaginationAndSort) ([]models.HabitActionSummary, error)
	GetHabitActionById(id string) (*models.HabitAction, error)
	CreateHabitAction(habitAction *models.HabitAction) error
	UpdateHabitAction(id string, habitAction *models.HabitAction) error
	DeleteHabitAction(id string) error
	ValidateEntry(habitAction *models.HabitAction) error
}

type habitActionService struct {
	repo      repositories.HabitActionRepository
	habitRepo repositories.HabitRepository
	userRepo  repositories.UserRepository
}

func NewHabitActionService(repo repositories.HabitActionRepository, habitRepo repositories.HabitRepository) HabitActionService {
	return &habitActionService{repo: repo, habitRepo: habitRepo}
}

func (s *habitActionService) GetHabitActions(query models.HabitActionQuery, pagination *structs.PaginationAndSort) ([]models.HabitAction, error) {
	return s.repo.GetHabitActions(&query, pagination)
}

func (s *habitActionService) GetHabitActionSummary(query models.HabitActionQuery, pagination *structs.PaginationAndSort) ([]models.HabitActionSummary, error) {
	return s.repo.GetHabitActionSummary(&query, pagination)
}

func (s *habitActionService) GetHabitActionById(id string) (*models.HabitAction, error) {
	return s.repo.GetHabitAction(id)
}

func (s *habitActionService) CreateHabitAction(habitAction *models.HabitAction) error {
	if err := s.ValidateEntry(habitAction); err != nil {
		return err
	}
	return s.repo.CreateHabitAction(habitAction)
}

func (s *habitActionService) UpdateHabitAction(id string, habitAction *models.HabitAction) error {
	if err := s.ValidateEntry(habitAction); err != nil {
		return err
	}
	return s.repo.UpdateHabitAction(id, habitAction)
}

func (s *habitActionService) DeleteHabitAction(id string) error {
	return s.repo.DeleteHabitAction(id)
}

func (s *habitActionService) ValidateEntry(habitAction *models.HabitAction) error {
	if s.isForCreate(habitAction) || habitAction.HabitID != "" {
		_, err := s.habitRepo.GetHabit(habitAction.HabitID)
		if err != nil {
			return err
		}
	}
	if s.isForCreate(habitAction) || habitAction.UserID != "" {
		_, err := s.userRepo.GetUser(habitAction.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *habitActionService) isForCreate(habitAction *models.HabitAction) bool {
	return habitAction.ID == ""
}
