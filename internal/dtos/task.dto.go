package dtos

import (
	"encoding/json"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"
	"time"
)

type SubtaskDTO struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TaskDTO struct {
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Title       string       `json:"title"`
	Notes       string       `json:"notes"`
	Icon        string       `json:"icon"`
	Color       string       `json:"color"`
	ArchivedAt  *time.Time   `json:"archived_at,omitempty"`
	Status      string       `json:"status"`
	DueDate     *string      `json:"due_date,omitempty"`
	DueTime     *string      `json:"due_time,omitempty"`
	Priority    string       `json:"priority"`
	Subtasks    []SubtaskDTO `json:"subtasks"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	UserID      string       `json:"user_id"`
	User        *UserDTO     `json:"user,omitempty"`
}

func ParseTask(task *models.Task) *TaskDTO {
	if task == nil {
		return nil
	}
	dto := &TaskDTO{
		ID:          task.ID,
		Type:        "task",
		Title:       task.Title,
		Notes:       task.Notes,
		Icon:        task.Icon,
		Color:       task.Color,
		ArchivedAt:  task.ArchivedAt,
		Status:      string(task.Status),
		DueTime:     task.DueTime,
		Priority:    task.Priority,
		Subtasks:    decodeSubtasks(task.Subtasks),
		CompletedAt: task.CompletedAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		UserID:      task.UserID,
		User:        ParseUser(task.User),
	}
	if dto.Status == "" {
		dto.Status = string(models.TaskStatusTodo)
	}
	if dto.Priority == "" {
		dto.Priority = "none"
	}
	if task.DueDate != nil && !task.DueDate.IsZero() {
		formatted := task.DueDate.Format("2006-01-02")
		dto.DueDate = &formatted
	}
	return dto
}

func ParseTasks(tasks []models.Task) []TaskDTO {
	result := make([]TaskDTO, 0, len(tasks))
	for i := range tasks {
		result = append(result, *ParseTask(&tasks[i]))
	}
	return result
}

type TaskCreateDTO struct {
	Title       string       `json:"title" binding:"required"`
	Notes       string       `json:"notes"`
	Icon        string       `json:"icon"`
	Color       string       `json:"color"`
	Status      string       `json:"status"`
	DueDate     string       `json:"due_date"`
	DueTime     *string      `json:"due_time"`
	Priority    string       `json:"priority"`
	Subtasks    []SubtaskDTO `json:"subtasks"`
}

func (t *TaskCreateDTO) ToTask() (*models.Task, error) {
	task := &models.Task{
		Title:    t.Title,
		Notes:    t.Notes,
		Icon:     t.Icon,
		Color:    t.Color,
		Status:   models.TaskStatus(t.Status),
		DueTime:  t.DueTime,
		Priority: t.Priority,
	}
	if task.Status == "" {
		task.Status = models.TaskStatusTodo
	}
	if task.Priority == "" {
		task.Priority = "none"
	}
	if t.DueDate != "" {
		dd, err := helpers.ParseDate(t.DueDate)
		if err != nil {
			return nil, err
		}
		task.DueDate = &dd
	}
	task.Subtasks = encodeSubtasks(t.Subtasks)
	return task, nil
}

type TaskUpdateDTO struct {
	Title       *string       `json:"title"`
	Notes       *string       `json:"notes"`
	Icon        *string       `json:"icon"`
	Color       *string       `json:"color"`
	Status      *string       `json:"status"`
	DueDate     *string       `json:"due_date"`
	DueTime     *string       `json:"due_time"`
	Priority    *string       `json:"priority"`
	Subtasks    *[]SubtaskDTO `json:"subtasks"`
	Archived    *bool         `json:"archived"`
}

func (t *TaskUpdateDTO) ToTask() (*models.Task, error) {
	task := &models.Task{}
	if t.Title != nil {
		task.Title = *t.Title
	}
	if t.Notes != nil {
		task.Notes = *t.Notes
	}
	if t.Icon != nil {
		task.Icon = *t.Icon
	}
	if t.Color != nil {
		task.Color = *t.Color
	}
	if t.Status != nil {
		task.Status = models.TaskStatus(*t.Status)
		if *t.Status == string(models.TaskStatusDone) {
			now := time.Now()
			task.CompletedAt = &now
		}
	}
	if t.DueDate != nil {
		dd, err := helpers.ParseDate(*t.DueDate)
		if err != nil {
			return nil, err
		}
		task.DueDate = &dd
	}
	if t.DueTime != nil {
		task.DueTime = t.DueTime
	}
	if t.Priority != nil {
		task.Priority = *t.Priority
	}
	if t.Subtasks != nil {
		task.Subtasks = encodeSubtasks(*t.Subtasks)
	}
	if t.Archived != nil {
		if *t.Archived {
			now := time.Now()
			task.ArchivedAt = &now
		} else {
			zero := time.Time{}
			task.ArchivedAt = &zero
		}
	}
	return task, nil
}

func encodeSubtasks(subtasks []SubtaskDTO) string {
	if subtasks == nil {
		subtasks = []SubtaskDTO{}
	}
	encoded, err := json.Marshal(subtasks)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}

func decodeSubtasks(raw string) []SubtaskDTO {
	subtasks := []SubtaskDTO{}
	if raw == "" {
		return subtasks
	}
	if err := json.Unmarshal([]byte(raw), &subtasks); err != nil {
		return []SubtaskDTO{}
	}
	return subtasks
}
