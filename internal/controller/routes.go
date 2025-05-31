package controller

import (
	"backend-mobAppRest/internal/controller/auth"
	"backend-mobAppRest/internal/controller/cart"
	"backend-mobAppRest/internal/controller/product"
	"backend-mobAppRest/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

func RouteAPI(route *gin.Engine, services *service.Service) {
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	apiRoute := route.Group("/api")
	apiRoute.Use(LoggerMiddleware())

	cartRoute := apiRoute.Group("/cart")
	cartController := cart.NewCartController(services.CartService)
	cartRoute.Use(AuthMiddleware())

	cartRoute.GET("/get", cartController.GetCart)
	cartRoute.POST("/plus", cartController.Plus)
	cartRoute.POST("/minus", cartController.Minus)

	catalogRoute := apiRoute.Group("/catalog")
	catalogController := product.NewCatalogController(services.ProductService)
	//catalogRoute.Use(cartController.AuthMiddleware())

	catalogRoute.GET("/categories", catalogController.GetCategories)
	catalogRoute.GET("/catalog", catalogController.GetCatalog)
	catalogRoute.GET("/product/:productID", catalogController.GetProductById)

	authRoute := apiRoute.Group("/auth")
	authController := auth.NewAuthController(services.AuthService)

	authRoute.POST("/signin", authController.SignIn)
	authRoute.POST("/signup", authController.SignUp)

	route.GET("/storage/:filename", func(c *gin.Context) {

		filename := c.Param("filename")

		if strings.Contains(filename, "..") {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Полный путь до файла
		filePath := filepath.Join("storage", filename)

		// Отправляем файл клиенту
		c.File(filePath)
	})
}
