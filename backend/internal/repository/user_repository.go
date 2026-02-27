package repository

import (
	"github.com/crimsonn/zm_server/internal/dto"
	"github.com/crimsonn/zm_server/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers() ([]models.User, error) {
	users := []models.User{}
	err := r.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) CreateOrUpdateUser(user *dto.CreateOrUpdateUserDTO) error {
	_, err := r.db.NamedExec(`
	INSERT INTO users (name, fivem, license) VALUES (:name, :fivem, :license)
	ON CONFLICT (fivem) DO UPDATE SET name = EXCLUDED.name, license = EXCLUDED.license, updated_at = NOW(), last_login = NOW()
	`, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByFivemIdentifier(identifier string) (*models.User, error) {
	user := models.User{}
	err := r.db.Get(&user, "SELECT * FROM users WHERE fivem = $1", identifier)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
