package cart

import (
	"backend-mobAppRest/internal/service"
	"backend-mobAppRest/internal/service/customServiceError"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ACartController struct {
	CartService service.CartService
}

func NewCartController(cartService service.CartService) *ACartController {
	return &ACartController{CartService: cartService}
}

func (ca *ACartController) GetCart(c *gin.Context) {
	token, _ := c.Get("token")

	carts, err := ca.CartService.GetCarts(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrCartsNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, carts)
}

func (ca *ACartController) Plus(c *gin.Context) {
	token, _ := c.Get("token")
	var json ActionRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	carts, err := ca.CartService.Plus(token.(string), json.ProductID)
	if err != nil {
		if errors.Is(err, customServiceError.ErrCartsNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, carts)
}

func (ca *ACartController) Minus(c *gin.Context) {
	token, _ := c.Get("token")
	var json ActionRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	carts, err := ca.CartService.Minus(token.(string), json.ProductID)
	if err != nil {
		if errors.Is(err, customServiceError.ErrCartsNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, carts)
}
