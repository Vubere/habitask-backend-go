package database

import (
	"fmt"
	"habitask-backend-go/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Habit{},
		&models.HabitCompletion{},
		&models.Task{},
		&models.Notification{},
		&models.DeviceToken{},
	)
	if err != nil {
		fmt.Println("DB Migration failed: ", err)
		panic(err)
	}
	fmt.Println("DB Migration successful!")
}
