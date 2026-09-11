package main

import (
	"habitask-backend-go/internal/controllers"
	"habitask-backend-go/internal/database"
	"habitask-backend-go/internal/middleware"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.Connect()
	database.Migrate(db)

	// Repositories
	userRepository := repositories.NewUserRepository(db)

	// Services
	userService := services.NewUserService(userRepository)

	// Controllers
	userController := controllers.NewUserController(userService)

	// Routes
	router := gin.Default()
	v1 := router.Group("/api/v1")

	// - Public
	// - - User
	v1.Group("/users")
	{
		v1.POST("/signup", userController.SignUp)
		v1.POST("/login", userController.Login)
	}

	// - Private
	v1.Use(middleware.AuthenticationMiddleware())
	{
		// - Users
		v1.Group("/users")
		{
			v1.GET("/me", userController.Me)
			v1.GET("/:id", userController.GetUserById)
			v1.GET("", userController.GetUsers)
			v1.PUT("/:id", userController.UpdateUser)
			v1.DELETE("/:id", userController.DeleteUser)
		}
	}
}
