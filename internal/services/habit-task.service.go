package services

import (
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/pkg/lib/structs"
)

type HabitTaskService interface {
	GetHabitTaskDisplayDetails(filter dtos.HabitTaskFilterDTO, pagination *structs.PaginationAndSort, search string) ([]dtos.HabitTaskDTO, error)
}

type habitTaskService struct {
	repo repositories.HabitTaskRepository
}

func NewHabitTaskService(repo repositories.HabitTaskRepository) HabitTaskService {
	return &habitTaskService{
		repo: repo,
	}
}
func (s *habitTaskService) GetHabitTaskDisplayDetails(filter dtos.HabitTaskFilterDTO, pagination *structs.PaginationAndSort, search string) ([]dtos.HabitTaskDTO, error) {
	return s.repo.GetHabitTaskDisplayDetails(filter, pagination, search)
}
