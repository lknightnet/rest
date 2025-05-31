package repository

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/repository/customRepositoryError"
	"backend-mobAppRest/pkg/database"
	"errors"
	"gorm.io/gorm"
)

type cartRepository struct {
	DB *database.PostgreSQL
}

func (c *cartRepository) GetCartsByAccessToken(accessToken string) ([]model.Cart, error) {
	var carts []model.Cart
	err := c.DB.DB.Where("access_token = ?", accessToken).Find(&carts).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrCartsNotFound
		}
		return nil, err
	}

	return carts, nil
}

func (c *cartRepository) CreateCart(cart *model.Cart) error {
	if err := c.DB.DB.Create(cart).Error; err != nil {
		return err
	}
	return nil
}

func (c *cartRepository) PlusCart(accessToken string, productID int) error {
	var item model.Cart
	err := c.DB.DB.Where("access_token = ? AND product_id = ?", accessToken, productID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customRepositoryError.ErrCartNotFound
		}
		return err
	}
	item.Quantity++
	return c.DB.DB.Save(&item).Error
}

func (c *cartRepository) MinusCart(accessToken string, productID int) (*model.Cart, error) {
	var item model.Cart
	err := c.DB.DB.Where("access_token = ? AND product_id = ?", accessToken, productID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrCartNotFound
		}
		return nil, err
	}

	// товар найден — увеличиваем количество
	item.Quantity--
	err = c.DB.DB.Save(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrCartNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (c *cartRepository) RemoveCart(accessToken string, productID int) error {
	var cart model.Cart
	err := c.DB.DB.Where("access_token = ? AND product_id = ?", accessToken, productID).Delete(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customRepositoryError.ErrCartNotFound
		}
		return err
	}
	return nil
}

func (c *cartRepository) ClearCarts(accessToken string) error {
	var cart model.Cart
	err := c.DB.DB.Where("access_token = ?", accessToken).Delete(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customRepositoryError.ErrCartNotFound
		}
		return err
	}
	return nil
}

func newCartRepository(db *database.PostgreSQL) *cartRepository {
	return &cartRepository{
		DB: db,
	}
}
