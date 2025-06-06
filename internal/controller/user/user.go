package user

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/service"
	"backend-mobAppRest/internal/service/customServiceError"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService}
}

func (u *UserController) GetUserByAccessToken(c *gin.Context) {
	token, _ := c.Get("token")

	user, err := u.UserService.GetUserByAccessToken(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrCartsNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *UserController) ChangeInformation(c *gin.Context) {
	token, _ := c.Get("token")
	var json ChangeRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := u.UserService.ChangeInformation(token.(string), &model.User{
		Name:     json.Name,
		Phone:    json.Phone,
		Password: json.Password,
		Address:  json.Address,
	})
	if err != nil {
		if errors.Is(err, customServiceError.ErrCartsNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, OKResponse{Status: true})
}
