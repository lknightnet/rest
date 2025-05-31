package product

import (
	"backend-mobAppRest/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CatalogController struct {
	ProductService service.ProductService
}

func NewCatalogController(productService service.ProductService) *CatalogController {
	return &CatalogController{
		ProductService: productService,
	}
}

func (p *CatalogController) GetCategories(c *gin.Context) {
	categories, err := p.ProductService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (p *CatalogController) GetCatalog(c *gin.Context) {
	catalog, err := p.ProductService.GetCatalog()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, catalog)
}

func (p *CatalogController) GetProductById(c *gin.Context) {
	productID := c.Param("productID")
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	catalog, err := p.ProductService.GetProductById(productIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, catalog)
}
