package service

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/repository"
	"time"
)

type CartService interface {
	GetCarts(token string) (*model.ViewCart, error)
	Plus(token string, productID int) (*model.ViewCart, error)
	Minus(token string, productID int) (*model.ViewCart, error)
	Clear(token string) (*model.ViewCart, error)
}

type ProductService interface {
	GetCategories() ([]model.ViewCategoryList, error)
	GetCatalog() ([]model.ViewCategoryWithProductList, error)
	GetProductById(productID int) (*model.Product, error)
}

type AuthService interface {
	SignUp(username, phone, password string) (*model.Tokens, error)
	SignIn(phone, password string) (*model.Tokens, error)
}

type UserService interface {
	GetUserByAccessToken(token string) (*model.ViewUser, error)
	ChangeInformation(token string, user *model.User) error
}

type Service struct {
	CartService    CartService
	ProductService ProductService
	AuthService    AuthService
	UserService    UserService
}

type DependenciesService struct {
	SignKey       []byte
	AuthSignature string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration

	AuthRepository    repository.AuthRepository
	TokenRepository   repository.TokenRepository
	UserRepository    repository.UserRepository
	CatalogRepository repository.CatalogRepository
	CartRepository    repository.CartRepository
}

func NewService(deps *DependenciesService) *Service {
	return &Service{
		AuthService: newAuthService(deps.SignKey, deps.AccessExpiry, deps.RefreshExpiry,
			deps.AuthSignature, deps.AuthRepository, deps.TokenRepository, deps.UserRepository),
		ProductService: newCatalogService(deps.CatalogRepository),
		CartService:    newCartService(deps.CartRepository, deps.CatalogRepository),
		UserService:    newUserService(deps.UserRepository),
	}
}
