package repository

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/pkg/database"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) (int, error)
	CreateOrderProductList(orderProductList *model.OrderProductList) error
	GetListOrders(userID int) ([]model.Order, error)
	GetOrderById(orderId int, userId int) (*model.Order, error)
	GetOrderProductListById(orderId int) ([]model.OrderProductList, error)
}

type AuthRepository interface {
	CreateUser(user *model.User) (int, error)
}
type UserRepository interface {
	GetUserByToken(token string) (*model.User, error)
	GetUserByPhone(phone string) (*model.User, error)
	GetUserByID(userID int) (*model.User, error) // noot use
	ChangeInformation(token string, user *model.User) error
}

type CartRepository interface {
	GetCartsByAccessToken(accessToken string) ([]model.Cart, error)
	CreateCart(cart *model.Cart) error
	PlusCart(accessToken string, productID int) error
	MinusCart(accessToken string, productID int) (*model.Cart, error)
	RemoveCart(accessToken string, productID int) error
	ClearCarts(accessToken string) error
}

type CatalogRepository interface {
	GetCategories() ([]model.Category, error)
	GetProductsByCategoryID(categoryID int) ([]model.Product, error)
	GetProductById(productID int) (*model.Product, error)
}

type TokenRepository interface {
	CreateTokens(accessToken *model.AccessToken, refreshToken *model.RefreshToken) error
	GetRefreshToken(refreshToken string) (*model.RefreshToken, error) // noot use
	GetAccessToken(accessToken string) (*model.AccessToken, error)    // noot use
}

type Repository struct {
	AuthRepository    AuthRepository
	UserRepository    UserRepository
	TokenRepository   TokenRepository
	CatalogRepository CatalogRepository
	CartRepository    CartRepository
	OrderRepository   OrderRepository
}

func NewRepositories(db *database.PostgreSQL) *Repository {
	return &Repository{
		AuthRepository:    newAutRepository(db),
		TokenRepository:   newTokenRepository(db),
		UserRepository:    newUserRepository(db),
		CatalogRepository: newCatalogRepository(db),
		CartRepository:    newCartRepository(db),
		OrderRepository:   newOrderRepository(db),
	}
}
