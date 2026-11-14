package repositories

import (
	"habitask-backend-go/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceTokenRepository interface {
	Upsert(token *models.DeviceToken) error
	GetByUser(userID string) ([]models.DeviceToken, error)
	DeleteByToken(token string) error
}

type deviceTokenRepository struct {
	db *gorm.DB
}

func NewDeviceTokenRepository(db *gorm.DB) DeviceTokenRepository {
	return &deviceTokenRepository{db: db}
}

func (r *deviceTokenRepository) Upsert(token *models.DeviceToken) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "token"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "platform", "updated_at"}),
	}).Create(token).Error
}

func (r *deviceTokenRepository) GetByUser(userID string) ([]models.DeviceToken, error) {
	tokens := []models.DeviceToken{}
	err := r.db.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

func (r *deviceTokenRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.DeviceToken{}).Error
}
