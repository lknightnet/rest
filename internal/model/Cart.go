package model

type Cart struct {
	ID          int
	ProductID   int
	Price       float64
	Quantity    int
	AccessToken string
}

func (c *Cart) GetProductID() int {
	return c.ProductID
}

type ViewCart struct {
	TotalPrice float64 `json:"total_price"`
	Product    []ViewProductCartList
}
