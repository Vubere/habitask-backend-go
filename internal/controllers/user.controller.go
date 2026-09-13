package controllers

import (
	"errors"
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
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	user := userInput.ToUser()
	err := c.userService.SignUp(user)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(201, gin.H{"message": "profile created successfully", "success": true})
}

func (c *userController) Login(ctx *gin.Context) {
	var userInput dtos.UserLoginDTO
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	token, user, err := c.userService.Login(userInput.Email, userInput.Password)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"token": token, "user": dtos.ParseUser(user), "success": true, "message": "user logged in successfully"})
}

func (c *userController) Me(ctx *gin.Context) {
	user, err := c.userService.GetUserById(ctx.GetString("userId"))
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.Set("userId", user.ID)
	ctx.JSON(200, gin.H{"user": dtos.ParseUser(user), "success": true, "message": "user retrieved successfully"})
}

func (c *userController) GetUserById(ctx *gin.Context) {
	user, err := c.userService.GetUserById(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"user": dtos.ParseUser(user), "success": true, "message": "user retrieved successfully"})
}

func (c *userController) GetUsers(ctx *gin.Context) {
	pagination, query, err := c.getUserPaginationAndQuery(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	users, err := c.userService.GetUsers(*query, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"data": users, "pagination": pagination, "success": true, "message": "users retrieved successfully"})
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	var userInput dtos.UserUpdateDTO
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	userToUpdate := userInput.ToUser()
	userToUpdate.ID = userId
	err := c.userService.UpdateUser(userId, userToUpdate)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "user updated successfully", "success": true})
}

func (c *userController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	err := c.userService.DeleteUser(userId)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error(), "success": false})
		return
	}
	ctx.JSON(200, gin.H{"message": "user deleted successfully", "success": true})
}

func (c *userController) getUserPaginationAndQuery(ctx *gin.Context) (*structs.PaginationAndSort, *models.UserQuery, error) {
	paginationAndSortPassed, exists := ctx.Get("paginationAndSort")
	if !exists {
		return nil, nil, errors.New("invalid pagination")
	}
	queryPassed, exists := ctx.Get("query")
	if !exists {
		return nil, nil, errors.New("invalid query")
	}
	query, ok := queryPassed.(*models.UserQuery)
	if !ok {
		return nil, nil, errors.New("invalid query")
	}
	paginationAndSort, ok := paginationAndSortPassed.(*structs.PaginationAndSort)
	if !ok {
		return nil, nil, errors.New("invalid pagination")
	}
	return paginationAndSort, query, nil
}
