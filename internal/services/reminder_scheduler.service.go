package services

import (
	"bytes"
	"encoding/json"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

const expoPushEndpoint = "https://exp.host/--/api/v2/push/send"

type ReminderSchedulerService interface {
	Start()
}

type reminderSchedulerService struct {
	db               *gorm.DB
	notificationRepo repositories.NotificationRepository
	deviceTokenRepo  repositories.DeviceTokenRepository
}

func NewReminderSchedulerService(
	db *gorm.DB,
	notificationRepo repositories.NotificationRepository,
	deviceTokenRepo repositories.DeviceTokenRepository,
) ReminderSchedulerService {
	return &reminderSchedulerService{
		db:               db,
		notificationRepo: notificationRepo,
		deviceTokenRepo:  deviceTokenRepo,
	}
}

func (s *reminderSchedulerService) Start() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		s.runOnce()
		for range ticker.C {
			s.runOnce()
		}
	}()
}

func (s *reminderSchedulerService) runOnce() {
	now := time.Now()
	hhmm := now.Format("15:04")
	today := now.Format("2006-01-02")

	habits := []models.Habit{}
	if err := s.db.Where("reminder_time = ? AND archived_at IS NULL", hhmm).Find(&habits).Error; err != nil {
		log.Println("reminder scheduler: habit query failed:", err)
	}
	for i := range habits {
		s.notify(habits[i].UserID, habits[i].Title, "Time for your habit", "habit", habits[i].ID)
	}

	tasks := []models.Task{}
	if err := s.db.
		Where("due_time = ? AND status = ? AND archived_at IS NULL AND (due_date IS NULL OR due_date = ?)", hhmm, models.TaskStatusTodo, today).
		Find(&tasks).Error; err != nil {
		log.Println("reminder scheduler: task query failed:", err)
	}
	for i := range tasks {
		s.notify(tasks[i].UserID, tasks[i].Title, "Task due now", "task", tasks[i].ID)
	}
}

func (s *reminderSchedulerService) notify(userID, title, body, source, sourceID string) {
	notification := &models.Notification{
		Title:       title,
		Description: body,
		Category:    "reminder",
		Date:        time.Now(),
		Source:      source,
		SourceID:    sourceID,
		UserID:      userID,
	}
	if err := s.notificationRepo.CreateNotification(notification); err != nil {
		log.Println("reminder scheduler: failed to store notification:", err)
	}

	tokens, err := s.deviceTokenRepo.GetByUser(userID)
	if err != nil || len(tokens) == 0 {
		return
	}
	pushTokens := make([]string, 0, len(tokens))
	for i := range tokens {
		pushTokens = append(pushTokens, tokens[i].Token)
	}
	sendExpoPush(pushTokens, title, body, map[string]interface{}{"source": source, "id": sourceID})
}

func sendExpoPush(tokens []string, title, body string, data map[string]interface{}) {
	if len(tokens) == 0 {
		return
	}
	messages := make([]map[string]interface{}, 0, len(tokens))
	for _, token := range tokens {
		messages = append(messages, map[string]interface{}{
			"to":    token,
			"title": title,
			"body":  body,
			"data":  data,
		})
	}
	payload, err := json.Marshal(messages)
	if err != nil {
		return
	}
	request, err := http.NewRequest(http.MethodPost, expoPushEndpoint, bytes.NewReader(payload))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		log.Println("reminder scheduler: push request failed:", err)
		return
	}
	defer response.Body.Close()
}
