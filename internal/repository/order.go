package repository

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/pkg/database"
)

type orderRepository struct {
	DB *database.PostgreSQL
}

func (o *orderRepository) CreateOrder(order *model.Order) (int, error) {
	if err := o.DB.DB.Create(order).Error; err != nil {
		return 0, err
	}
	return order.ID, nil
}

func (o *orderRepository) CreateOrderProductList(orderProductList *model.OrderProductList) error {
	if err := o.DB.DB.Create(orderProductList).Error; err != nil {
		return err
	}
	return nil
}

func (o *orderRepository) GetListOrders(userID int) ([]model.Order, error) {
	var orders []model.Order

	err := o.DB.DB.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderRepository) GetOrderById(orderId int, userId int) (*model.Order, error) {
	var order model.Order

	err := o.DB.DB.
		Where("id = ? AND user_id = ?", orderId, userId).
		First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *orderRepository) GetOrderProductListById(orderId int) ([]model.OrderProductList, error) {
	var orders []model.OrderProductList

	err := o.DB.DB.Where("order_id = ?", orderId).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func newOrderRepository(db *database.PostgreSQL) *orderRepository {
	return &orderRepository{
		DB: db,
	}
}
