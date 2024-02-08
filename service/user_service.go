package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/petruskuswandi/bwastartup.git/models"
	"github.com/petruskuswandi/bwastartup.git/repository"
	"github.com/petruskuswandi/bwastartup.git/request"
	"golang.org/x/crypto/bcrypt"
)

type ServiceUser interface {
	RegisterUser(req request.RegisterUserInput) (models.User, error)
	Login(input request.LoginInput) (models.User, error)
	IsEmailAvailable(input request.CheckEmailInput) (bool, error)
	SaveAvatar(ID, fileLocation string) (models.User, error)
	GetUserByID(ID string) (models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(input request.FormUpdateUserInput) (models.User, error)
}

type serviceUser struct {
	repo repository.UserRepository
}

func NewServiceUser(repo repository.UserRepository) *serviceUser {
	return &serviceUser{repo: repo}
}
func (s *serviceUser) RegisterUser(req request.RegisterUserInput) (models.User, error) {
	user := models.User{}

	userID, err := uuid.NewRandom()
	if err != nil {
		return user, err
	}
	user.ID = userID.String()
	user.Name = req.Name
	user.Email = req.Email
	user.Occupation = req.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repo.SaveUser(user)
	if err != nil {
		return newUser, err
	}

	return user, nil
}

func (s *serviceUser) Login(input request.LoginInput) (models.User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repo.FindByEmailUser(email)
	if err != nil {
		return user, err
	}

	if user.ID == "" {
		return user, errors.New("no user found on that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *serviceUser) IsEmailAvailable(input request.CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repo.FindByEmailUser(email)
	if err != nil {
		return false, err
	}

	if user.ID == "" {
		return true, nil
	}

	return false, nil
}

func (s *serviceUser) SaveAvatar(ID, fileLocation string) (models.User, error) {
	user, err := s.repo.FindByIDUser(ID)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updateUser, err := s.repo.UpdateUser(user)
	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

func (s *serviceUser) GetUserByID(ID string) (models.User, error) {
	user, err := s.repo.FindByIDUser(ID)
	if err != nil {
		return user, err
	}

	if user.ID == "" {
		return user, errors.New("no user found on with that ID")
	}

	return user, nil
}

func (s *serviceUser) GetAllUsers() ([]models.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return users, err
	}

	return users, nil
}

func (s *serviceUser) UpdateUser(input request.FormUpdateUserInput) (models.User, error) {
	user, err := s.repo.FindByIDUser(input.ID)
	if err != nil {
		return user, err
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation

	updatedUser, err := s.repo.UpdateUser(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}