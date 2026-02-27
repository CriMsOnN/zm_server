package user

import (
	"github.com/crimsonn/zm_server/internal/dto"
	"github.com/crimsonn/zm_server/internal/models"
	"github.com/crimsonn/zm_server/internal/repository"
)

type UserService struct {
	repository  *repository.UserRepository
	onlineUsers map[string]models.User
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.repository.GetUsers()
}

func (s *UserService) CreateOrUpdateUser(user *dto.CreateOrUpdateUserDTO) error {
	return s.repository.CreateOrUpdateUser(user)
}
