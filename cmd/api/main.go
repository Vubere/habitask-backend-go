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
	notificationRepository := repositories.NewNotificationRepository(db)
	taskRepository := repositories.NewTaskRepository(db)
	habitRepository := repositories.NewHabitRepository(db)
	habitCompletionRepository := repositories.NewHabitCompletionRepository(db)
	habitTaskRepository := repositories.NewHabitTaskRepository(db)
	syncRepository := repositories.NewSyncRepository(db)
	deviceTokenRepository := repositories.NewDeviceTokenRepository(db)

	// Services
	userService := services.NewUserService(userRepository)
	notificationService := services.NewNotificationService(notificationRepository, userRepository)
	taskService := services.NewTaskService(taskRepository, userRepository)
	habitService := services.NewHabitService(habitRepository, userRepository)
	habitCompletionService := services.NewHabitCompletionService(habitCompletionRepository, habitRepository, userRepository)
	habitTaskService := services.NewHabitTaskService(habitTaskRepository)
	syncService := services.NewSyncService(syncRepository)
	deviceTokenService := services.NewDeviceTokenService(deviceTokenRepository)
	reminderScheduler := services.NewReminderSchedulerService(db, notificationRepository, deviceTokenRepository)

	// Controllers
	userController := controllers.NewUserController(userService)
	notificationController := controllers.NewNotificationController(notificationService)
	taskController := controllers.NewTaskController(taskService)
	habitController := controllers.NewHabitController(habitService)
	habitCompletionController := controllers.NewHabitCompletionController(habitCompletionService)
	habitTaskController := controllers.NewHabitTaskController(habitTaskService)
	syncController := controllers.NewSyncController(syncService)
	deviceTokenController := controllers.NewDeviceTokenController(deviceTokenService)

	// Background reminder scheduler
	reminderScheduler.Start()

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
	public := v1.Group("")
	// - - User
	private := v1.Group("")

	public.Group("")
	{
		public.POST("/users/signup", userController.SignUp)
		public.POST("/users/login", userController.Login)
	}

	// - Private
	private.Use(middleware.AuthenticationMiddleware())
	{
		// - Users
		private.GET("/users/me", userController.Me)
		private.GET("/users/:id", userController.GetUserById)
		private.GET("/users", middleware.UserMiddleware(), userController.GetUsers)
		private.PUT("/users/:id", userController.UpdateUser)
		private.DELETE("/users/:id", userController.DeleteUser)

		// - Notification
		private.GET("/notifications", middleware.NotificationMiddleware(), notificationController.GetNotifications)
		private.GET("/notifications/:id", notificationController.GetNotificationById)
		private.GET("/notifications/summary", middleware.NotificationMiddleware(), notificationController.GetNotificationSummary)
		private.PUT("/notifications/mark-as-read-/:id", notificationController.MarkNotificationAsRead)
		private.PUT("/notifications/mark-all-as-read", notificationController.MarkAllNotificationsAsRead)

		// - Tasks
		private.GET("/tasks", middleware.TaskMiddleware(), taskController.GetTasks)
		private.GET("/tasks/:id", taskController.GetTaskById)
		private.GET("/tasks/summary", middleware.TaskMiddleware(), taskController.GetTaskSummary)
		private.POST("/tasks", taskController.CreateTask)
		private.PUT("/tasks/:id", taskController.UpdateTask)
		private.DELETE("/tasks/:id", taskController.DeleteTask)

		// - Habits
		private.GET("/habits", middleware.HabitMiddleware(), habitController.GetHabits)
		private.GET("/habits/summary", middleware.HabitMiddleware(), habitController.GetHabitSummary)
		private.GET("/habits/:id/stats", habitController.GetHabitStats)
		private.GET("/habits/:id", habitController.GetHabitById)
		private.POST("/habits", habitController.CreateHabit)
		private.POST("/habits/:id/completions", habitCompletionController.LogCompletion)
		private.PUT("/habits/:id", habitController.UpdateHabit)
		private.DELETE("/habits/:id", habitController.DeleteHabit)

		// - Habit Completions
		private.GET("/habit-completions", middleware.HabitCompletionMiddleware(), habitCompletionController.GetHabitCompletions)
		private.GET("/habit-completions/summary", middleware.HabitCompletionMiddleware(), habitCompletionController.GetHabitCompletionSummary)
		private.GET("/habit-completions/:id", habitCompletionController.GetHabitCompletionById)
		private.POST("/habit-completions", habitCompletionController.CreateHabitCompletion)
		private.PUT("/habit-completions/:id", habitCompletionController.UpdateHabitCompletion)
		private.DELETE("/habit-completions/:id", habitCompletionController.DeleteHabitCompletion)

		// - Habit Tasks
		private.GET("/habit-task", middleware.HabitTaskMiddleware(), habitTaskController.GetHabitTaskDisplayDetails)

		// - Sync
		private.GET("/sync", syncController.GetChanges)

		// - Devices (push tokens)
		private.POST("/devices", deviceTokenController.Register)
		private.DELETE("/devices/:token", deviceTokenController.Unregister)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
