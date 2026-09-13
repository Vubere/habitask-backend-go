package main

import (
	"fmt"
	"habitask-backend-go/internal/controllers"
	"habitask-backend-go/internal/database"
	"habitask-backend-go/internal/middleware"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/internal/services"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := database.Connect()
	database.Migrate(db)

	// Repositories
	userRepository := repositories.NewUserRepository(db)
	reminderRepository := repositories.NewReminderRepository(db)
	notificationRepository := repositories.NewNotificationRepository(db)
	taskRepository := repositories.NewTaskRepository(db)
	taskStepRepository := repositories.NewTaskStepsRepository(db)
	habitRepository := repositories.NewHabitRepository(db)
	habitActionRepository := repositories.NewHabitActionRepository(db)
	habitOccurenceRepository := repositories.NewHabitOccurenceRepository(db)

	// Services
	userService := services.NewUserService(userRepository)
	reminderService := services.NewReminderService(reminderRepository, userRepository)
	notificationService := services.NewNotificationService(notificationRepository, userRepository)
	taskService := services.NewTaskService(taskRepository, userRepository)
	taskStepService := services.NewTaskStepService(taskStepRepository, userRepository, taskRepository)
	habitService := services.NewHabitService(habitRepository, userRepository)
	habitActionService := services.NewHabitActionService(habitActionRepository, habitRepository)
	habitOccurenceService := services.NewHabitOccurenceService(habitOccurenceRepository, habitRepository)

	// Controllers
	userController := controllers.NewUserController(userService)
	reminderController := controllers.NewReminderController(reminderService)
	notificationController := controllers.NewNotificationController(notificationService)
	taskController := controllers.NewTaskController(taskService)
	taskStepController := controllers.NewTaskStepsController(taskStepService)
	habitController := controllers.NewHabitController(habitService)
	habitActionController := controllers.NewHabitActionController(habitActionService)
	habitOccurenceController := controllers.NewHabitOccurencesController(habitOccurenceService)

	// Routes
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    404,
			"message": "The requested resource was not found",
		})
	})

	v1 := router.Group("/api/v1")

	// - Public
	// - - User
	public := v1.Group("")
	private := v1.Group("")

	public.Group("")
	{
		v1.POST("/users/signup", userController.SignUp)
		v1.POST("/users/login", userController.Login)
	}

	// - Private
	private.Use(middleware.AuthenticationMiddleware())
	{
		// - Users
		v1.GET("/users/me", userController.Me)
		v1.GET("/users/:id", userController.GetUserById)
		v1.GET("/users", middleware.UserMiddleware(), userController.GetUsers)
		v1.PUT("/users/:id", userController.UpdateUser)
		v1.DELETE("/users/:id", userController.DeleteUser)

		// - Notifications
		v1.GET("/notifications", notificationController.GetNotifications)
		v1.GET("/notifications/:id", notificationController.GetNotificationById)
		v1.GET("/notifications/summary", notificationController.GetNotificationSummary)
		v1.PUT("/notifications/mark-as-read-/:id", notificationController.MarkNotificationAsRead)
		v1.PUT("/notifications/mark-all-as-read", notificationController.MarkAllNotificationsAsRead)

		// - Tasks
		v1.GET("/tasks", taskController.GetTasks)
		v1.GET("/tasks/:id", taskController.GetTaskById)
		v1.GET("/tasks/summary", taskController.GetTaskSummary)
		v1.POST("/tasks", taskController.CreateTask)
		v1.PUT("/tasks/:id", taskController.UpdateTask)
		v1.DELETE("/tasks/:id", taskController.DeleteTask)

		// - Task Steps
		v1.GET("/tasks-steps", taskStepController.GetTaskSteps)
		v1.GET("/tasks-steps/:id", taskStepController.GetTaskStepById)
		v1.GET("/tasks-steps/summary", taskStepController.GetTaskStepSummary)
		v1.POST("/tasks-steps", taskStepController.CreateTaskStep)
		v1.PUT("/tasks-steps/:id", taskStepController.UpdateTaskStep)
		v1.DELETE("/tasks-steps/:id", taskStepController.DeleteTaskStep)

		// - Habits
		v1.GET("/habits", habitController.GetHabits)
		v1.GET("/habits/:id", habitController.GetHabitById)
		v1.GET("/habits/summary", habitController.GetHabitSummary)
		v1.POST("/habits", habitController.CreateHabit)
		v1.PUT("/habits/:id", habitController.UpdateHabit)
		v1.DELETE("/habits/:id", habitController.DeleteHabit)

		// - Habit Actions
		v1.GET("/habits-actions", habitActionController.GetHabitActions)
		v1.GET("/habits-actions/:id", habitActionController.GetHabitActionById)
		v1.GET("/habits-actions/summary", habitActionController.GetHabitActionSummary)
		v1.POST("/habits-actions", habitActionController.CreateHabitAction)
		v1.PUT("/habits-actions/:id", habitActionController.UpdateHabitAction)
		v1.DELETE("/habits-actions/:id", habitActionController.DeleteHabitAction)

		// - Habit Occurences
		v1.GET("/habits-occurences", habitOccurenceController.GetHabitOccurences)
		v1.GET("/habits-occurences/:id", habitOccurenceController.GetHabitOccurenceById)
		v1.GET("/habits-occurences/summary", habitOccurenceController.GetHabitOccurenceSummary)
		v1.POST("/habits-occurences", habitOccurenceController.CreateHabitOccurence)
		v1.PUT("/habits-occurences/:id", habitOccurenceController.UpdateHabitOccurence)

		// - Reminders
		v1.GET("/reminders", reminderController.GetReminders)
		v1.GET("/reminders/:reminder_id", reminderController.GetReminderById)
		v1.POST("/reminders", reminderController.CreateReminder)
		v1.PUT("/reminders/:reminder_id", reminderController.UpdateReminder)
		v1.DELETE("/reminders/:reminder_id", reminderController.DeleteReminder)
	}
	// not found

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
