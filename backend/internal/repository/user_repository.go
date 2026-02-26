package repository

import (
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
