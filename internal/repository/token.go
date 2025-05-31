package repository

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/pkg/database"
	"errors"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *database.PostgreSQL
}

func newTokenRepository(db *database.PostgreSQL) *tokenRepository {
	return &tokenRepository{db: db}
}
func (t *tokenRepository) CreateTokens(accessToken *model.AccessToken, refreshToken *model.RefreshToken) error {
	tx := t.db.DB.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	if err := t.db.DB.Create(accessToken).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := t.db.DB.Create(refreshToken).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (t *tokenRepository) GetRefreshToken(refreshToken string) (*model.RefreshToken, error) {
	var refresh model.RefreshToken

	if err := t.db.DB.Where("token = ?", refreshToken).First(&refresh).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found")
		}
		return nil, err
	}
	return &refresh, nil
}

func (t *tokenRepository) GetAccessToken(accessToken string) (*model.AccessToken, error) {
	var access model.AccessToken

	if err := t.db.DB.Where("token = ?", accessToken).First(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("access token not found")
		}
		return nil, err
	}
	return &access, nil
}
