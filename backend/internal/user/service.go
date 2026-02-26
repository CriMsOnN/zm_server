package user

import (
	"github.com/crimsonn/zm_server/internal/models"
	"github.com/crimsonn/zm_server/internal/repository"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.repository.GetUsers()
}
