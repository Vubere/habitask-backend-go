package services

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/pkg/lib/structs"
)

type NotificationService interface {
	GetNotifications(query models.NotificationQuery, pagination *structs.PaginationAndSort) ([]models.Notification, error)
	GetNotificationSummary(query models.NotificationQuery, pagination *structs.PaginationAndSort) ([]models.NotificationSummary, error)
	GetNotificationById(id string) (*models.Notification, error)
	CreateNotification(notification *models.Notification) error
	UpdateNotification(id string, notification *models.Notification) error
	MarkNotificationAsRead(id string) error
	MarkAllNotificationsAsRead(userId string) error
	DeleteNotification(id string) error
	ValidateEntry(notification *models.Notification) error
}

type notificationService struct {
	repo     repositories.NotificationRepository
	userRepo repositories.UserRepository
}

func NewNotificationService(repo repositories.NotificationRepository, userRepo repositories.UserRepository) NotificationService {
	return &notificationService{repo: repo, userRepo: userRepo}
}

func (s *notificationService) GetNotifications(query models.NotificationQuery, pagination *structs.PaginationAndSort) ([]models.Notification, error) {
	return s.repo.GetNotifications(&query, pagination)
}

func (s *notificationService) GetNotificationSummary(query models.NotificationQuery, pagination *structs.PaginationAndSort) ([]models.NotificationSummary, error) {
	return s.repo.GetNotificationSummary(&query, pagination)
}

func (s *notificationService) GetNotificationById(id string) (*models.Notification, error) {
	return s.repo.GetNotification(id)
}

func (s *notificationService) CreateNotification(notification *models.Notification) error {
	if err := s.ValidateEntry(notification); err != nil {
		return err
	}
	return s.repo.CreateNotification(notification)
}

func (s *notificationService) UpdateNotification(id string, notification *models.Notification) error {
	if err := s.ValidateEntry(notification); err != nil {
		return err
	}
	return s.repo.UpdateNotification(id, notification)
}

func (s *notificationService) MarkNotificationAsRead(id string) error {
	return s.repo.MarkNotificationAsRead(id)
}

func (s *notificationService) MarkAllNotificationsAsRead(userId string) error {
	return s.repo.MarkAllNotificationsAsRead(userId)
}

func (s *notificationService) DeleteNotification(id string) error {
	return s.repo.DeleteNotification(id)
}

func (s *notificationService) ValidateEntry(notification *models.Notification) error {
	if s.isForCreate(notification) || notification.UserID != "" {
		_, err := s.userRepo.GetUser(notification.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *notificationService) isForCreate(notification *models.Notification) bool {
	return notification.ID == ""
}
