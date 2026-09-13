package services

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/pkg/lib/structs"
)

type HabitOccurenceService interface {
	GetHabitOccurences(query models.HabitOccurenceQuery, pagination *structs.PaginationAndSort) ([]models.HabitOccurence, error)
	GetHabitOccurenceSummary(query models.HabitOccurenceQuery, pagination *structs.PaginationAndSort) ([]models.HabitOccurenceSummary, error)
	GetHabitOccurenceById(id string) (*models.HabitOccurence, error)
	CreateHabitOccurence(habitOccurence *models.HabitOccurence) error
	UpdateHabitOccurence(id string, habitOccurence *models.HabitOccurence) error
	DeleteHabitOccurence(id string) error
}

type habitOccurenceService struct {
	repo      repositories.HabitOccurenceRepository
	habitRepo repositories.HabitRepository
	userRepo  repositories.UserRepository
}

func NewHabitOccurenceService(repo repositories.HabitOccurenceRepository, habitRepo repositories.HabitRepository) HabitOccurenceService {
	return &habitOccurenceService{repo: repo, habitRepo: habitRepo}
}

func (s *habitOccurenceService) GetHabitOccurences(query models.HabitOccurenceQuery, pagination *structs.PaginationAndSort) ([]models.HabitOccurence, error) {
	return s.repo.GetHabitOccurences(&query, pagination)
}

func (s *habitOccurenceService) GetHabitOccurenceSummary(query models.HabitOccurenceQuery, pagination *structs.PaginationAndSort) ([]models.HabitOccurenceSummary, error) {
	return s.repo.GetHabitOccurenceSummary(&query, pagination)
}

func (s *habitOccurenceService) GetHabitOccurenceById(id string) (*models.HabitOccurence, error) {
	return s.repo.GetHabitOccurence(id)
}

func (s *habitOccurenceService) CreateHabitOccurence(habitOccurence *models.HabitOccurence) error {
	if err := s.ValidateEntry(habitOccurence); err != nil {
		return err
	}
	return s.repo.CreateHabitOccurence(habitOccurence)
}

func (s *habitOccurenceService) UpdateHabitOccurence(id string, habitOccurence *models.HabitOccurence) error {
	if err := s.ValidateEntry(habitOccurence); err != nil {
		return err
	}
	return s.repo.UpdateHabitOccurence(habitOccurence)
}

func (s *habitOccurenceService) DeleteHabitOccurence(id string) error {
	return s.repo.DeleteHabitOccurence(id)
}

func (s *habitOccurenceService) ValidateEntry(habitOccurence *models.HabitOccurence) error {
	if s.isForCreate(habitOccurence) || habitOccurence.HabitID != "" {
		_, err := s.habitRepo.GetHabit(habitOccurence.HabitID)
		if err != nil {
			return err
		}
	}
	if s.isForCreate(habitOccurence) || habitOccurence.UserID != "" {
		_, err := s.userRepo.GetUser(habitOccurence.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *habitOccurenceService) isForCreate(habitOccurence *models.HabitOccurence) bool {
	return habitOccurence.ID == ""
}
