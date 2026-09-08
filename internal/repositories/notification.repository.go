package repositories

import (
	"fmt"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"
	"habitask-backend-go/pkg/lib/database/scopes"
	"habitask-backend-go/pkg/lib/structs"
	"strings"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	GetNotifications(filter *models.NotificationQuery, pagination *structs.PaginationAndSort) ([]*models.Notification, error)
	GetNotification(id string) (*models.Notification, error)
	GetNotificationSummary(filter *models.NotificationQuery, pagination *structs.PaginationAndSort) ([]models.NotificationSummary, error)
	CreateNotification(notification *models.Notification) error
	UpdateNotification(notification *models.Notification) error
	DeleteNotification(id string) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) GetNotifications(filter *models.NotificationQuery, pagination *structs.PaginationAndSort) ([]*models.Notification, error) {
	notifications := []*models.Notification{}
	query := r.db.Model(&models.Notification{}).Where(filter.Notification)
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.DateLte != nil {
		query = query.Where("date <= ?", *filter.DateLte)
	}
	if filter.DateGte != nil {
		query = query.Where("date >= ?", *filter.DateGte)
	}
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) GetNotification(id string) (*models.Notification, error) {
	notification := &models.Notification{}
	err := r.db.First(notification, id).Error
	return notification, err
}

func (r *notificationRepository) GetNotificationSummary(filter *models.NotificationQuery, pagination *structs.PaginationAndSort) ([]models.NotificationSummary, error) {
	var notifications []models.NotificationSummary
	groupExpr, err := helpers.GetSummaryExpression(helpers.SummaryGroup{
		GroupBy:         filter.GroupBy,
		DateGroup:       filter.DateGroup,
		Field:           filter.GroupBy,
		AllowedGrouping: notificationSummaryGroups,
	})
	if err != nil {
		return nil, err
	}
	query := r.db.Model(&models.Notification{}).Where(filter.Notification).Select(
		fmt.Sprintf(`
		%s as label,
		COUNT(notifications.id) as count,
		MAX(notifications.date) as last_notification_date,
		MAX(notifications.is_read) as is_read
		`, groupExpr),
	)
	if strings.HasPrefix(groupExpr, "users.") {
		query = query.Joins("JOIN users ON notifications.user_id = users.id")
	}
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.DateLte != nil {
		query = query.Where("date <= ?", *filter.DateLte)
	}
	if filter.DateGte != nil {
		query = query.Where("date >= ?", *filter.DateGte)
	}
	err = query.Offset(pagination.GetOffset()).Limit(pagination.PerPage).Group(groupExpr).Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *notificationRepository) CreateNotification(notification *models.Notification) error {
	err := r.db.Create(notification).Error
	return err
}

func (r *notificationRepository) UpdateNotification(notification *models.Notification) error {
	err := r.db.Where("id = ?", notification.ID).Updates(notification).Error
	return err
}

func (r *notificationRepository) DeleteNotification(id string) error {
	notification := &models.Notification{}
	err := r.db.Where("id = ?", id).Delete(notification).Error
	return err
}

var notificationSummaryGroups = map[string]string{
	"created_at": "notifications.created_at",
	"date":       "notifications.date",
	"title":      "notifications.title",
	"category":   "notifications.category",
	"is_read":    "notifications.is_read",
	"user_id":    "notifications.user_id",
	"username":   "users.username",
	"source":     "notifications.source",
	"source_id":  "notifications.source_id",
}
