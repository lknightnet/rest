package repository

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/repository/customRepositoryError"
	"backend-mobAppRest/pkg/database"
	"errors"
	"gorm.io/gorm"
)

type catalogRepository struct {
	DB *database.PostgreSQL
}

func (c *catalogRepository) GetProductById(productID int) (*model.Product, error) {
	var product model.Product

	if err := c.DB.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (c *catalogRepository) GetCategories() ([]model.Category, error) {
	var categories []model.Category

	err := c.DB.DB.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *catalogRepository) GetProductsByCategoryID(categoryID int) ([]model.Product, error) {
	var products []model.Product
	err := c.DB.DB.Where("category_id = ?", categoryID).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func newCatalogRepository(db *database.PostgreSQL) *catalogRepository {
	return &catalogRepository{
		DB: db,
	}
}
