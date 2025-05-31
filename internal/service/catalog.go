package service

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/repository"
	"backend-mobAppRest/internal/service/customServiceError"
	"backend-mobAppRest/pkg/tg"
	"log/slog"
	"sort"
)

type catalogService struct {
	CatalogRepository repository.CatalogRepository
}

func (c *catalogService) GetCategories() ([]model.ViewCategoryList, error) {
	var catalogs []model.ViewCategoryList

	categories, err := c.CatalogRepository.GetCategories()
	if err != nil {
		go tg.SendError(err.Error(), "/api/catalog/categories")
		slog.Debug("error get categories", err)
		return nil, customServiceError.ErrUnknown
	}

	for _, category := range categories {
		catalogs = append(catalogs, model.ViewCategoryList{
			ID:       category.ID,
			Name:     category.Name,
			Sort:     category.Sort,
			ImageUrl: "/storage/" + category.ImageUrl,
		})
	}

	sort.SliceIsSorted(catalogs, func(i, j int) bool {
		return catalogs[i].Sort < catalogs[j].Sort
	})
	return catalogs, nil
}

func (c *catalogService) GetCatalog() ([]model.ViewCategoryWithProductList, error) {
	var catalogs []model.ViewCategoryWithProductList

	categories, err := c.CatalogRepository.GetCategories()
	if err != nil {
		go tg.SendError(err.Error(), "/api/catalog/catalog")
		slog.Debug("error get categories", err)
		return nil, customServiceError.ErrUnknown
	}

	for _, category := range categories {

		products, err := c.CatalogRepository.GetProductsByCategoryID(category.ID)
		if err != nil {
			go tg.SendError(err.Error(), "/api/catalog/catalog")
			slog.Debug("error get products", err)
			return nil, customServiceError.ErrUnknown
		}

		productList := make([]model.ViewProductList, 0, len(products))
		for _, p := range products {
			productList = append(productList, model.ViewProductList{
				ID:         p.ID,
				Name:       p.Name,
				CategoryID: category.ID,
				Image:      p.Image,
				Weight:     p.Weight,
				Price:      p.Price,
			})
		}

		catalogs = append(catalogs, model.ViewCategoryWithProductList{
			ID:          category.ID,
			Name:        category.Name,
			Sort:        category.Sort,
			ProductList: productList,
		})
	}
	return catalogs, nil
}

func (c *catalogService) GetProductById(productID int) (*model.Product, error) {
	product, err := c.CatalogRepository.GetProductById(productID)
	if err != nil {
		go tg.SendError(err.Error(), "/api/catalog/product")
		slog.Debug("error get product", err)
		return nil, customServiceError.ErrUnknown
	}
	return product, nil
}

func newCatalogService(catalogRepository repository.CatalogRepository) *catalogService {
	return &catalogService{
		CatalogRepository: catalogRepository,
	}
}
