package services

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/pkg/lib/structs"
)

type HabitService interface {
	GetHabits(query models.HabitQuery, pagination *structs.PaginationAndSort) ([]models.Habit, error)
	GetHabitSummary(query models.HabitQuery, pagination *structs.PaginationAndSort) ([]models.HabitSummary, error)
	GetHabitById(id string) (*models.Habit, error)
	CreateHabit(habit *models.Habit) error
	UpdateHabit(id string, habit *models.Habit) error
	DeleteHabit(id string) error
	ValidateEntry(habit *models.Habit) error
}

type habitService struct {
	repo     repositories.HabitRepository
	userRepo repositories.UserRepository
}

func NewHabitService(repo repositories.HabitRepository, userRepo repositories.UserRepository) HabitService {
	return &habitService{repo: repo, userRepo: userRepo}
}

func (s *habitService) GetHabits(query models.HabitQuery, pagination *structs.PaginationAndSort) ([]models.Habit, error) {
	return s.repo.GetHabits(&query, pagination)
}

func (s *habitService) GetHabitSummary(query models.HabitQuery, pagination *structs.PaginationAndSort) ([]models.HabitSummary, error) {
	return s.repo.GetHabitSummary(&query, pagination)
}

func (s *habitService) GetHabitById(id string) (*models.Habit, error) {
	return s.repo.GetHabit(id)
}

func (s *habitService) CreateHabit(habit *models.Habit) error {
	if err := s.ValidateEntry(habit); err != nil {
		return err
	}
	return s.repo.CreateHabit(habit)
}

func (s *habitService) UpdateHabit(id string, habit *models.Habit) error {
	if err := s.ValidateEntry(habit); err != nil {
		return err
	}
	return s.repo.UpdateHabit(id, habit)
}

func (s *habitService) DeleteHabit(id string) error {
	return s.repo.DeleteHabit(id)
}

func (s *habitService) ValidateEntry(habit *models.Habit) error {
	if s.isForCreate(habit) || habit.UserID != "" {
		_, err := s.userRepo.GetUser(habit.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *habitService) isForCreate(habit *models.Habit) bool {
	return habit.ID == ""
}
