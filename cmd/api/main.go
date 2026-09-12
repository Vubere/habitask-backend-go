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

	// Services
	userService := services.NewUserService(userRepository)

	// Controllers
	userController := controllers.NewUserController(userService)

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
		v1.GET("/users", userController.GetUsers)
		v1.PUT("/users/:id", userController.UpdateUser)
		v1.DELETE("/users/:id", userController.DeleteUser)
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
