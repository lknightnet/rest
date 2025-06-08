package model

type Order struct {
	ID                      int
	TotalPrice              float64
	InstrumentationQuantity int
	UserID                  int
	Address                 string
	UserPhone               string
	City                    string
	Status                  string
	IsDelivery              bool
	PaymentMethod           bool
	Bonuses                 int
	Comment                 string
}

type ViewOrderByID struct {
	ID                      int     `json:"id"`
	TotalPrice              float64 `json:"total_price"`
	InstrumentationQuantity int     `json:"instrumentation_quantity"`
	Address                 string  `json:"address"`
	UserPhone               string  `json:"user_phone"`
	City                    string  `json:"city"`
	Status                  string  `json:"status"`
	IsDelivery              bool    `json:"is_delivery"`
	PaymentMethod           bool    `json:"payment_method"`
	Bonuses                 int     `json:"bonuses"`
	Comment                 string  `json:"comment"`
}

type ViewOrderByIDList struct {
	OrderInfo        ViewOrderByID          `json:"order_info"`
	OrderProductList []ViewOrderProductList `json:"order_product_list"`
}

type ViewOrdersByID []ViewOrderByIDList

type ViewOrder struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type ViewOrderList []ViewOrder
