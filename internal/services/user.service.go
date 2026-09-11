package services

import (
	"errors"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/internal/repositories"
	"habitask-backend-go/pkg/appbcrypt"
	"habitask-backend-go/pkg/appjwt"
	"habitask-backend-go/pkg/lib/structs"

	"gorm.io/gorm"
)

type UserService interface {
	SignUp(user *models.User) error
	Login(usernameOrEmail string, password string) (string, *models.User, error)
	GetUserById(userId string) (*models.User, error)
	GetUsers(query models.UserQuery, pagination *structs.PaginationAndSort) ([]models.User, error)
	UpdateUser(userId string, user *models.User) error
	DeleteUser(userId string) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) SignUp(user *models.User) error {
	rec, err := s.userRepository.GetUserByEmail(user.Email)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return err
		}
	}
	if rec.ID != "" {
		return errors.New("email already exists")
	}
	rec, err = s.userRepository.GetUserByUsername(user.Username)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return err
		}
	}
	if rec != nil {
		return errors.New("username already exists")
	}
	user.Password, err = appbcrypt.HashPassword(user.Password)
	if err != nil {
		return err
	}
	return s.userRepository.CreateUser(user)
}

func (s *userService) Login(usernameOrEmail string, password string) (string, *models.User, error) {
	user, err := s.userRepository.GetByUsernameOrEmail(usernameOrEmail)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return "", nil, errors.New("invalid username or password")
		}
		return "", nil, err
	}
	if appbcrypt.VerifyPassword(password, user.Password) {
		return "", nil, errors.New("invalid username or password")
	}
	token, err := appjwt.CreateToken(appjwt.JwtUserInfo{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *userService) GetUserById(userId string) (*models.User, error) {
	return s.userRepository.GetUser(userId)
}

func (s *userService) GetUsers(query models.UserQuery, pagination *structs.PaginationAndSort) ([]models.User, error) {
	return s.userRepository.GetUsers(&query, pagination)
}

func (s *userService) UpdateUser(userId string, user *models.User) error {
	return s.userRepository.UpdateUser(user)
}

func (s *userService) DeleteUser(userId string) error {
	return s.userRepository.DeleteUser(userId)
}
