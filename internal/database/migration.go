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
		&models.HabitAction{},
		&models.Reminder{},
		&models.Task{},
		&models.TaskSteps{},
		&models.Notification{},
		&models.HabitOccurence{},
	)
	if err != nil {
		fmt.Println("DB Migration failed: ", err)
		panic(err)
	}
}
