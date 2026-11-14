package services

import (
	"time"

	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/repositories"
)

type SyncService interface {
	GetChanges(userID string, since time.Time) (*dtos.SyncDTO, error)
}

type syncService struct {
	repo repositories.SyncRepository
}

func NewSyncService(repo repositories.SyncRepository) SyncService {
	return &syncService{repo: repo}
}

func (s *syncService) GetChanges(userID string, since time.Time) (*dtos.SyncDTO, error) {
	return s.repo.GetChanges(userID, since)
}
