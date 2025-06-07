package service

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/repository"
	"backend-mobAppRest/internal/repository/customRepositoryError"
	"backend-mobAppRest/internal/service/customServiceError"
	"backend-mobAppRest/pkg/tg"
	"errors"
	"log/slog"
)

type cartService struct {
	CatalogRepository repository.CatalogRepository
	CartRepository    repository.CartRepository
}

func (c *cartService) GetCarts(token string) (*model.ViewCart, error) {
	carts, err := c.CartRepository.GetCartsByAccessToken(token)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrCartsNotFound) {
			return nil, customServiceError.ErrCartsNotFound
		}

		go tg.SendError(err.Error(), "/api/cart/get")
		slog.Debug("error get cart by access token", err, token)
		return nil, customServiceError.ErrUnknown
	}

	viewCart := &model.ViewCart{}
	var viewProductCartList = make([]model.ViewProductCartList, 0)

	for _, cart := range carts {
		viewCart.TotalPrice = viewCart.TotalPrice + float64(cart.Quantity)*cart.Price

		product, err := c.CatalogRepository.GetProductById(cart.ProductID)
		if err != nil {
			go tg.SendError(err.Error(), "/api/cart/get")
			slog.Debug("error get product by id", err, cart.ProductID)
			return nil, customServiceError.ErrUnknown
		}

		viewProductCartList = append(viewProductCartList, model.ViewProductCartList{
			ID:       cart.ProductID,
			Name:     product.Name,
			Image:    "/storage/" + product.Image,
			Weight:   product.Weight,
			Quantity: cart.Quantity,
			Price:    product.Price,
		})
	}
	viewCart.Product = viewProductCartList

	return viewCart, nil
}

func (c *cartService) Plus(token string, productID int) (*model.ViewCart, error) {

	err := c.CartRepository.PlusCart(token, productID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrCartNotFound) {

			product, err := c.CatalogRepository.GetProductById(productID)
			if err != nil {
				if errors.Is(err, customRepositoryError.ErrProductNotFound) {
					go tg.SendError(err.Error(), "/api/cart/plus")
					slog.Debug("error get product by id", err, productID)
					return nil, customRepositoryError.ErrProductNotFound
				}
				go tg.SendError(err.Error(), "/api/cart/plus")
				slog.Debug("error get product by id", err, productID)
				return nil, customServiceError.ErrUnknown
			}

			cart := &model.Cart{
				ProductID:   productID,
				Price:       product.Price,
				Quantity:    1,
				AccessToken: token,
			}

			err = c.CartRepository.CreateCart(cart)
			if err != nil {
				go tg.SendError(err.Error(), "/api/cart/plus")
				slog.Debug("error create cart", err, cart)
				return nil, customServiceError.ErrUnknown
			}
			return c.GetCarts(token)
		}
		go tg.SendError(err.Error(), "/api/cart/plus")
		slog.Debug("error plus cart", err, token, productID)
		return nil, customServiceError.ErrUnknown
	}

	return c.GetCarts(token)
}

func (c *cartService) Minus(token string, productID int) (*model.ViewCart, error) {
	cart, err := c.CartRepository.MinusCart(token, productID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrCartNotFound) {
			go tg.SendError(err.Error(), "/api/cart/minus")
			slog.Debug("error minus cart", err, token, productID)
			return nil, customRepositoryError.ErrCartNotFound
		}
		go tg.SendError(err.Error(), "/api/cart/minus")
		slog.Debug("error minus cart", err, token, productID)
		return nil, customServiceError.ErrUnknown
	}

	if cart.Quantity < 1 {
		err = c.CartRepository.RemoveCart(token, productID)
		if err != nil {
			if errors.Is(err, customRepositoryError.ErrCartNotFound) {
				go tg.SendError(err.Error(), "/api/cart/minus")
				slog.Debug("cart not found during remove", err, token, productID)
				return nil, customRepositoryError.ErrCartNotFound
			}
			go tg.SendError(err.Error(), "/api/cart/minus")
			slog.Debug("error remove cart", err, token, productID)
			return nil, customServiceError.ErrUnknown
		}
		return c.GetCarts(token)
	}

	return c.GetCarts(token)
}

func (c *cartService) Clear(token string) (*model.ViewCart, error) {
	//TODO implement me
	panic("implement me")
}

func newCartService(cartRepository repository.CartRepository, CatalogRepository repository.CatalogRepository) *cartService {
	return &cartService{CartRepository: cartRepository,
		CatalogRepository: CatalogRepository}
}
