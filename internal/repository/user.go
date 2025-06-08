package repository

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/repository/customRepositoryError"
	"backend-mobAppRest/pkg/database"
	"errors"
	"gorm.io/gorm"
)

type userRepository struct {
	DB *database.PostgreSQL
}

func (u *userRepository) ChangeInformation(token string, user *model.User) error {
	var accessToken model.AccessToken
	existingUser := &model.User{}
	err := u.DB.DB.Where("token = ?", token).First(&accessToken).Error
	if err != nil {
		return err
	}

	err = u.DB.DB.Where("id = ?", accessToken.Subject).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customRepositoryError.ErrUserNotFound
		}
		return err
	}

	// Создаем карту обновлений
	updates := make(map[string]interface{})

	if user.Name != "" {
		updates["name"] = user.Name
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}
	if user.Phone != "" {
		updates["phone"] = user.Phone
	}
	if user.Address != "" {
		updates["address"] = user.Address
	}
	if len(updates) == 0 {
		return nil
	}

	return u.DB.DB.Model(existingUser).Updates(updates).Error
}

func (u *userRepository) GetUserByToken(token string) (*model.User, error) {
	var user model.User
	var accessToken model.AccessToken
	err := u.DB.DB.Where("token = ?", token).First(&accessToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrUserNotFound
		}
		return nil, err
	}

	err = u.DB.DB.Where("id = ?", accessToken.Subject).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	err := u.DB.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrInvalidEmailOrPassword
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetUserByID(userID int) (*model.User, error) {
	var user model.User
	err := u.DB.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func newUserRepository(db *database.PostgreSQL) *userRepository {
	return &userRepository{
		DB: db,
	}
}
