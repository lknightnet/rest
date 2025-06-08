package controller

import (
	"backend-mobAppRest/internal/controller/auth"
	"backend-mobAppRest/internal/controller/cart"
	"backend-mobAppRest/internal/controller/order"
	"backend-mobAppRest/internal/controller/product"
	"backend-mobAppRest/internal/controller/user"
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
	RouteStorage(route)

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

	userRoute := apiRoute.Group("/user")
	userController := user.NewUserController(services.UserService)
	userRoute.Use(AuthMiddleware())

	userRoute.GET("/get", userController.GetUserByAccessToken)
	userRoute.POST("/change", userController.ChangeInformation)

	orderRoute := apiRoute.Group("/order")
	orderController := order.NewOrderController(services.OrderService)
	orderRoute.Use(AuthMiddleware())
	orderRoute.GET("/list", orderController.ListOrder)
	orderRoute.GET("/get/:id", orderController.OrderByID)
	orderRoute.POST("/create", orderController.Order)
}

func RouteStorage(route *gin.Engine) {
	storageRoute := route.Group("/storage")
	storageRoute.GET("/categories/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		if strings.Contains(filename, "..") {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Полный путь до файла
		filePath := filepath.Join("storage/categories", filename)

		// Отправляем файл клиенту
		c.File(filePath)
	})
	storageRoute.GET("/products/:filename", func(c *gin.Context) {

		filename := c.Param("filename")

		if strings.Contains(filename, "..") {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Полный путь до файла
		filePath := filepath.Join("storage/products", filename)

		// Отправляем файл клиенту
		c.File(filePath)
	})

	// Загрузка изображения категории
	storageRoute.POST("/categories", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить файл"})
			return
		}

		if strings.Contains(file.Filename, "..") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное имя файла"})
			return
		}

		savePath := filepath.Join("storage/categories", filepath.Base(file.Filename))
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Файл загружен", "path": "/storage/categories/" + file.Filename})
	})
	storageRoute.POST("/products", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить файл"})
			return
		}

		if strings.Contains(file.Filename, "..") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное имя файла"})
			return
		}

		savePath := filepath.Join("storage/products", filepath.Base(file.Filename))
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Файл загружен", "path": "/storage/products/" + file.Filename})
	})
}
