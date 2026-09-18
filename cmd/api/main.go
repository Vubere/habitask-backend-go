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
	habitTaskRepository := repositories.NewHabitTaskRepository(db)

	// Services
	userService := services.NewUserService(userRepository)
	reminderService := services.NewReminderService(reminderRepository, userRepository)
	notificationService := services.NewNotificationService(notificationRepository, userRepository)
	taskService := services.NewTaskService(taskRepository, userRepository)
	taskStepService := services.NewTaskStepService(taskStepRepository, userRepository, taskRepository)
	habitService := services.NewHabitService(habitRepository, userRepository)
	habitActionService := services.NewHabitActionService(habitActionRepository, habitRepository)
	habitOccurenceService := services.NewHabitOccurenceService(habitOccurenceRepository, habitRepository)
	habitTaskService := services.NewHabitTaskService(habitTaskRepository)

	// Controllers
	userController := controllers.NewUserController(userService)
	reminderController := controllers.NewReminderController(reminderService)
	notificationController := controllers.NewNotificationController(notificationService)
	taskController := controllers.NewTaskController(taskService)
	taskStepController := controllers.NewTaskStepsController(taskStepService)
	habitController := controllers.NewHabitController(habitService)
	habitActionController := controllers.NewHabitActionController(habitActionService)
	habitOccurenceController := controllers.NewHabitOccurencesController(habitOccurenceService)
	habitTaskController := controllers.NewHabitTaskController(habitTaskService)

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

		// - Notifications
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

		// - Task Steps
		private.GET("/tasks-steps", middleware.TaskStepMiddleware(), taskStepController.GetTaskSteps)
		private.GET("/tasks-steps/:id", taskStepController.GetTaskStepById)
		private.GET("/tasks-steps/summary", middleware.TaskStepMiddleware(), taskStepController.GetTaskStepSummary)
		private.POST("/tasks-steps", taskStepController.CreateTaskStep)
		private.PUT("/tasks-steps/:id", taskStepController.UpdateTaskStep)
		private.DELETE("/tasks-steps/:id", taskStepController.DeleteTaskStep)

		// - Habits
		private.GET("/habits", middleware.HabitMiddleware(), habitController.GetHabits)
		private.GET("/habits/:id", habitController.GetHabitById)
		private.GET("/habits/summary", middleware.HabitMiddleware(), habitController.GetHabitSummary)
		private.POST("/habits", habitController.CreateHabit)
		private.PUT("/habits/:id", habitController.UpdateHabit)
		private.DELETE("/habits/:id", habitController.DeleteHabit)

		// - Habit Actions
		private.GET("/habits-actions", middleware.HabitActionMiddleware(), habitActionController.GetHabitActions)
		private.GET("/habits-actions/:id", habitActionController.GetHabitActionById)
		private.GET("/habits-actions/summary", middleware.HabitActionMiddleware(), habitActionController.GetHabitActionSummary)
		private.POST("/habits-actions", habitActionController.CreateHabitAction)
		private.PUT("/habits-actions/:id", habitActionController.UpdateHabitAction)
		private.DELETE("/habits-actions/:id", habitActionController.DeleteHabitAction)

		// - Habit Occurences
		private.GET("/habits-occurences", middleware.HabitOccurenceMiddleware(), habitOccurenceController.GetHabitOccurences)
		private.GET("/habits-occurences/:id", habitOccurenceController.GetHabitOccurenceById)
		private.GET("/habits-occurences/summary", middleware.HabitOccurenceMiddleware(), habitOccurenceController.GetHabitOccurenceSummary)
		private.POST("/habits-occurences", habitOccurenceController.CreateHabitOccurence)
		private.PUT("/habits-occurences/:id", habitOccurenceController.UpdateHabitOccurence)

		// - Reminders
		private.GET("/reminders", middleware.ReminderMiddleware(), reminderController.GetReminders)
		private.GET("/reminders/:reminder_id", reminderController.GetReminderById)
		private.POST("/reminders", reminderController.CreateReminder)
		private.PUT("/reminders/:reminder_id", reminderController.UpdateReminder)
		private.DELETE("/reminders/:reminder_id", reminderController.DeleteReminder)

		// - Habit Tasks
		private.GET("/habit-task", middleware.HabitTaskMiddleware(), habitTaskController.GetHabitTaskDisplayDetails)
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
