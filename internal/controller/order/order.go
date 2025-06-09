package order

import (
	"backend-mobAppRest/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OrderController struct {
	OrderService service.OrderService
}

func NewOrderController(OrderService service.OrderService) *OrderController {
	return &OrderController{
		OrderService: OrderService,
	}
}

func (o *OrderController) Order(c *gin.Context) {
	token, _ := c.Get("token")

	var json OrderRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	id, err := o.OrderService.Order(token.(string), json.InstrumentationQuantity, json.IsDelivery, json.PaymentMethod, json.City, json.Bonuses, json.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (o *OrderController) ListOrder(c *gin.Context) {
	token, _ := c.Get("token")

	orderList, err := o.OrderService.ListOrder(token.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_list": orderList,
	})
}

func (o *OrderController) OrderByID(c *gin.Context) {
	token, _ := c.Get("token")

	orderID := c.Param("id")
	orderIDInt, err := strconv.Atoi(orderID)

	order, err := o.OrderService.OrderByID(token.(string), orderIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}
