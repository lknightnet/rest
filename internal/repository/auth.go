package repository

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/pkg/database"
)

type autRepository struct {
	DB *database.PostgreSQL
}

func (a *autRepository) CreateUser(user *model.User) (int, error) {
	if err := a.DB.DB.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func newAutRepository(db *database.PostgreSQL) *autRepository {
	return &autRepository{
		DB: db,
	}
}
