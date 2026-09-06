package repositories

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/lib/database/scopes"
	"habitask-backend-go/pkg/lib/structs"

	"gorm.io/gorm"
)

type ReminderRepository interface {
	GetReminders(filter *models.ReminderQuery, pagination *structs.PaginationAndSort) ([]*models.Reminder, error)
	GetReminder(id string) (*models.Reminder, error)
	CreateReminder(reminder *models.Reminder) error
	UpdateReminder(reminder *models.Reminder) error
	DeleteReminder(id string) error
}

type reminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepository{db: db}
}

func (r *reminderRepository) GetReminders(filter *models.ReminderQuery, pagination *structs.PaginationAndSort) ([]*models.Reminder, error) {
	reminders := []*models.Reminder{}
	query := r.db.Model(&models.Reminder{}).Where(filter.Reminder)
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.DateLte != nil {
		query = query.Where("date <= ?", *filter.DateLte)
	}
	if filter.DateGte != nil {
		query = query.Where("date >= ?", *filter.DateGte)
	}
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetReminder(id string) (*models.Reminder, error) {
	reminder := &models.Reminder{}
	err := r.db.First(reminder, id).Error
	return reminder, err
}

func (r *reminderRepository) CreateReminder(reminder *models.Reminder) error {
	err := r.db.Create(reminder).Error
	return err
}

func (r *reminderRepository) UpdateReminder(reminder *models.Reminder) error {
	err := r.db.Where("id = ?", reminder.ID).Updates(reminder).Error
	return err
}

func (r *reminderRepository) DeleteReminder(id string) error {
	reminder := &models.Reminder{}
	err := r.db.Where("id = ?", id).Delete(reminder).Error
	return err
}
