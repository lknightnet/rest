package service

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/repository"
	"backend-mobAppRest/internal/repository/customRepositoryError"
	"backend-mobAppRest/internal/service/customServiceError"
	"backend-mobAppRest/pkg/tg"
	"errors"
	"log"
	"log/slog"
)

type orderService struct {
	CartRepository    repository.CartRepository
	UserRepository    repository.UserRepository
	OrderRepository   repository.OrderRepository
	CatalogRepository repository.CatalogRepository
}

func (o *orderService) OrderByID(token string, orderId int) (model.ViewOrdersByID, error) {
	user, err := o.UserRepository.GetUserByToken(token)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrUserNotFound) {
			return nil, customServiceError.ErrUserNotFound
		}
		go tg.SendError(err.Error(), "/api/order/get/:id")
		slog.Debug("error get user by access token in database", err, token)
		return nil, customServiceError.ErrUnknown
	}

	order, err := o.OrderRepository.GetOrderById(orderId, user.ID)
	if err != nil {
		go tg.SendError(err.Error(), "/api/order/get/:id")
		slog.Debug("error get user by access token in database", err, token)
		return nil, customServiceError.ErrUnknown
	}

	orderProductList, err := o.OrderRepository.GetOrderProductListById(orderId)
	if err != nil {
		go tg.SendError(err.Error(), "/api/order/get/:id")
		slog.Debug("error get user by access token in database", err, token)
		return nil, customServiceError.ErrUnknown
	}

	viewOrderList := model.ViewOrdersByID{}
	orderProducts := []model.ViewOrderProductList{}

	for _, product := range orderProductList {

		productCatalog, err := o.CatalogRepository.GetProductById(product.ProductID)
		if err != nil {
			go tg.SendError(err.Error(), "/api/order/get/:id")
			slog.Debug("error get user by access token in database", err, token)
			return nil, customServiceError.ErrUnknown
		}

		orderProducts = append(orderProducts, model.ViewOrderProductList{
			Price:       product.Price,
			Quantity:    product.Quantity,
			ProductName: productCatalog.Name, // предположим, что ты подставишь тут настоящее имя
		})
	}

	viewOrderList = append(viewOrderList, model.ViewOrderByIDList{
		OrderInfo: model.ViewOrderByID{
			ID:                      order.ID,
			TotalPrice:              order.TotalPrice,
			InstrumentationQuantity: order.InstrumentationQuantity,
			Address:                 order.Address,
			UserPhone:               order.UserPhone,
			City:                    order.City,
			Status:                  order.Status,
			IsDelivery:              order.IsDelivery,
			PaymentMethod:           order.PaymentMethod,
			Bonuses:                 order.Bonuses,
			Comment:                 order.Comment,
		},
		OrderProductList: orderProducts,
	})

	return viewOrderList, nil

}

func (o *orderService) Order(token string, instrumentationQuantity int, isDelivery bool, paymentMethod string, City string, Bonuses int, Comment string) (int, error) {
	order := &model.Order{
		InstrumentationQuantity: instrumentationQuantity,
		City:                    City,
		Status:                  "Готовится",
		IsDelivery:              isDelivery,
		PaymentMethod:           paymentMethod,
		Bonuses:                 Bonuses,
		Comment:                 Comment,
	}

	user, err := o.UserRepository.GetUserByToken(token)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrUserNotFound) {
			return 0, customServiceError.ErrUserNotFound
		}
		go tg.SendError(err.Error(), "/api/order/create")
		slog.Debug("error get user by access token in database", err, token)
		return 0, customServiceError.ErrUnknown
	}

	order.UserID = user.ID
	order.Address = user.Address
	order.UserPhone = user.Phone

	carts, err := o.CartRepository.GetCartsByAccessToken(token)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrCartsNotFound) {
			return 0, customServiceError.ErrCartsNotFound
		}

		go tg.SendError(err.Error(), "/api/order/create")
		slog.Debug("error get cart by access token", err, token)
		return 0, customServiceError.ErrUnknown
	}

	totalPrice := 0.0
	for _, cart := range carts {
		totalPrice += float64(cart.Quantity) * cart.Price
	}

	order.TotalPrice = totalPrice

	log.Println(order)

	id, err := o.OrderRepository.CreateOrder(order)
	if err != nil {
		go tg.SendError(err.Error(), "/api/order/create")
		slog.Debug("error create order", err, order)
		return 0, customServiceError.ErrUnknown
	}

	for _, cart := range carts {
		orderProductList := model.ConstructCartToOrderProductList(&cart, id)
		err := o.OrderRepository.CreateOrderProductList(orderProductList)
		if err != nil {
			go tg.SendError(err.Error(), "/api/order/create")
			slog.Debug("error create order product list", err, order)
			return 0, customServiceError.ErrUnknown
		}
	}

	return id, nil
}

func (o *orderService) ListOrder(token string) (model.ViewOrderList, error) {
	user, err := o.UserRepository.GetUserByToken(token)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrUserNotFound) {
			return nil, customServiceError.ErrUserNotFound
		}
		go tg.SendError(err.Error(), "/api/order/list")
		slog.Debug("error get user by access token in database", err, token)
		return nil, customServiceError.ErrUnknown
	}

	log.Println(user.ID)

	orders, err := o.OrderRepository.GetListOrders(user.ID)
	if err != nil {
		go tg.SendError(err.Error(), "/api/order/list")
		slog.Debug("error get user by access token in database", err, token)
		return nil, customServiceError.ErrUnknown
	}

	log.Println(orders)

	viewOrderList := model.ViewOrderList{}
	for _, order := range orders {
		log.Println(order)
		viewOrderList = append(viewOrderList, model.ViewOrder{
			ID:     order.ID,
			Status: order.Status,
		})
	}
	return viewOrderList, nil
}

func newOrderService(cartRepository repository.CartRepository, userRepository repository.UserRepository, orderRepository repository.OrderRepository, catalogRepository repository.CatalogRepository) *orderService {
	return &orderService{
		CartRepository:    cartRepository,
		UserRepository:    userRepository,
		OrderRepository:   orderRepository,
		CatalogRepository: catalogRepository,
	}
}
