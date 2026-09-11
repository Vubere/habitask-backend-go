package controllers

import (
	"habitask-backend-go/internal/dtos"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/services"
	"habitask-backend-go/pkg/lib/structs"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	SignUp(ctx *gin.Context)
	Login(ctx *gin.Context)
	Me(ctx *gin.Context)
	GetUserById(ctx *gin.Context)
	GetUsers(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{userService: userService}
}

func (c *userController) SignUp(ctx *gin.Context) {
	var userInput dtos.UserCreateDTO
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := userInput.ToUser()
	err := c.userService.SignUp(user)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, gin.H{"message": "profile created successfully"})
}

func (c *userController) Login(ctx *gin.Context) {
	var userInput dtos.UserCreateDTO
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, user, err := c.userService.Login(userInput.Email, userInput.Password)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"token": token, "user": dtos.ParseUser(user)})
}

func (c *userController) Me(ctx *gin.Context) {
	user, err := c.userService.GetUserById(ctx.GetString("userId"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.Set("userId", user.ID)
	ctx.JSON(200, gin.H{"user": dtos.ParseUser(user)})
}

func (c *userController) GetUserById(ctx *gin.Context) {
	user, err := c.userService.GetUserById(ctx.Param("userId"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"user": dtos.ParseUser(user)})
}

func (c *userController) GetUsers(ctx *gin.Context) {
	query := models.UserQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	paginationPassed, exists := ctx.Get("pagination")
	if !exists {
		ctx.JSON(400, gin.H{"error": "invalid pagination"})
		return
	}
	pagination, ok := paginationPassed.(*structs.PaginationAndSort)
	if !ok {
		ctx.JSON(400, gin.H{"error": "invalid pagination"})
		return
	}
	users, err := c.userService.GetUsers(query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"users": users})
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	var userInput dtos.UserUpdateDTO
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userToUpdate := userInput.ToUser()
	userToUpdate.ID = userId
	err := c.userService.UpdateUser(userId, userToUpdate)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "user updated successfully"})
}

func (c *userController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	err := c.userService.DeleteUser(userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "user deleted successfully"})
}
