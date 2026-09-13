package services

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/pkg/lib/structs"
)

type ReminderService interface {
	GetReminders(query models.ReminderQuery, pagination *structs.PaginationAndSort) ([]models.Reminder, error)
	GetReminderById(id string) (*models.Reminder, error)
	CreateReminder(reminder *models.Reminder) error
	UpdateReminder(id string, reminder *models.Reminder) error
	DeleteReminder(id string) error
	ValidateEntry(reminder *models.Reminder) error
}

type reminderService struct {
	repo     repositories.ReminderRepository
	userRepo repositories.UserRepository
}

func NewReminderService(repo repositories.ReminderRepository, userRepo repositories.UserRepository) ReminderService {
	return &reminderService{repo: repo, userRepo: userRepo}
}

func (s *reminderService) GetReminders(query models.ReminderQuery, pagination *structs.PaginationAndSort) ([]models.Reminder, error) {
	return s.repo.GetReminders(&query, pagination)
}

func (s *reminderService) GetReminderById(id string) (*models.Reminder, error) {
	return s.repo.GetReminder(id)
}

func (s *reminderService) CreateReminder(reminder *models.Reminder) error {
	if err := s.ValidateEntry(reminder); err != nil {
		return err
	}
	return s.repo.CreateReminder(reminder)
}

func (s *reminderService) UpdateReminder(id string, reminder *models.Reminder) error {
	if err := s.ValidateEntry(reminder); err != nil {
		return err
	}
	return s.repo.UpdateReminder(id, reminder)
}

func (s *reminderService) DeleteReminder(id string) error {
	return s.repo.DeleteReminder(id)
}

func (s *reminderService) ValidateEntry(reminder *models.Reminder) error {
	if s.isForCreate(reminder) || reminder.UserID != "" {
		_, err := s.userRepo.GetUser(reminder.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *reminderService) isForCreate(reminder *models.Reminder) bool {
	return reminder.ID == ""
}
