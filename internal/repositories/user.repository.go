package repositories

import (
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/lib/database/scopes"
	"habitask-backend-go/pkg/lib/structs"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUsers(filter *models.UserQuery, pagination *structs.PaginationAndSort) ([]models.User, error)
	GetUser(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByToken(token string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetByUsernameOrEmail(usernameOrEmail string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUsers(filter *models.UserQuery, pagination *structs.PaginationAndSort) ([]models.User, error) {
	users := []models.User{}
	query := r.db.Model(&models.User{}).Where(filter.User)
	if filter.Search != "" {
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.Active != nil {
		query = query.Where("is_active = ?", *filter.Active)
	}
	if filter.Admin != nil {
		query = query.Where("is_admin = ?", *filter.Admin)
	}
	err := query.Scopes(scopes.ApplyPaginationAndSort(pagination)).Find(&users).Error
	return users, err
}

func (r *userRepository) GetUser(id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.First(user, id).Error
	return user, err
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("email = ?", email).First(user).Error
	return user, err
}

func (r *userRepository) GetUserByToken(token string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("token = ?", token).First(user).Error
	return user, err
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("username = ?", username).First(user).Error
	return user, err
}

func (r *userRepository) GetByUsernameOrEmail(usernameOrEmail string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("users.username = ? OR users.email = ?", usernameOrEmail, usernameOrEmail).First(user).Error
	return user, err
}

func (r *userRepository) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error
	return err
}

func (r *userRepository) UpdateUser(user *models.User) error {
	err := r.db.Save(user).Error
	return err
}

func (r *userRepository) DeleteUser(id string) error {
	user := &models.User{}
	err := r.db.Where("id = ?", id).Delete(user).Error
	return err
}
