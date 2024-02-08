package repository

import (
	"github.com/petruskuswandi/bwastartup.git/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	SaveUser(user models.User) (models.User, error)
	FindByEmailUser(email string) (models.User, error)
	FindByIDUser(ID string) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	FindAll() ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewRepositoryUser(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) SaveUser(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindByEmailUser(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindByIDUser(ID string) (models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", ID).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(user models.User) (models.User, error) {
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User

	err := r.db.Find(&users).Error
	if err != nil {
		return users, err
	}

	return users, nil
}
