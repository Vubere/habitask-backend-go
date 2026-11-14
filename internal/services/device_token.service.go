package services

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
)

type DeviceTokenService interface {
	Register(userID string, token *models.DeviceToken) error
	Unregister(token string) error
}

type deviceTokenService struct {
	repo repositories.DeviceTokenRepository
}

func NewDeviceTokenService(repo repositories.DeviceTokenRepository) DeviceTokenService {
	return &deviceTokenService{repo: repo}
}

func (s *deviceTokenService) Register(userID string, token *models.DeviceToken) error {
	token.UserID = userID
	return s.repo.Upsert(token)
}

func (s *deviceTokenService) Unregister(token string) error {
	return s.repo.DeleteByToken(token)
}
