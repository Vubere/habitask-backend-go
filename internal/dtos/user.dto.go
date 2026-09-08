package dtos

import "habitask-backend-go/internal/models"

type UserDTO struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Profession string `json:"profession"`
	Username   string `json:"user_name"`
	Bio        string `json:"bio"`
}

func ParseUser(user *models.User) *UserDTO {
	return &UserDTO{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Profession: user.Profession,
		Bio:        user.Bio,
		Username:   user.Username,
	}
}

type UserCreateDTO struct {
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Username   string `json:"user_name" binding:"required"`
	Profession string `json:"profession"`
	Bio        string `json:"bio"`
}

func (u *UserCreateDTO) ToUser() *models.User {
	user := &models.User{
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Email:      u.Email,
		Password:   u.Password,
		Username:   u.Username,
		Profession: u.Profession,
		Bio:        u.Bio,
		IsAdmin:    false,
		IsActive:   true,
	}
	return user
}

type UserUpdateDTO struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	Email      *string `json:"email"`
	Password   *string `json:"password"`
	Profession *string `json:"profession"`
	Bio        *string `json:"bio"`
	Username   *string `json:"user_name"`
}

func (u *UserUpdateDTO) ToUser() *models.User {
	user := &models.User{}
	if u.FirstName != nil {
		user.FirstName = *u.FirstName
	}
	if u.LastName != nil {
		user.LastName = *u.LastName
	}
	if u.Email != nil {
		user.Email = *u.Email
	}
	if u.Password != nil {
		user.Password = *u.Password
	}
	if u.Profession != nil {
		user.Profession = *u.Profession
	}
	if u.Bio != nil {
		user.Bio = *u.Bio
	}
	if u.Username != nil {
		user.Username = *u.Username
	}
	return user
}
