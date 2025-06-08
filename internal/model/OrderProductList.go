package model

type OrderProductList struct {
	ID        int
	Price     float64
	Quantity  int
	ProductID int
	OrderID   int
}

func ConstructCartToOrderProductList(cart *Cart, orderID int) *OrderProductList {
	return &OrderProductList{
		Price:     cart.Price,
		Quantity:  cart.Quantity,
		ProductID: cart.ProductID,
		OrderID:   orderID,
	}
}

type ViewOrderProductList struct {
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	ProductName string  `json:"product_name"`
}
