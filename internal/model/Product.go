package model

type Product struct {
	ID            int
	CategoryID    string
	Name          string
	Image         string
	Calorie       int
	Weight        int
	Price         float64
	Squirrels     int
	Fats          int
	Carbohydrates int
	Visibility    bool
}

type ViewProductList struct {
	ID         int
	Name       string
	CategoryID int
	Image      string
	Weight     int
	Price      float64
}

type ViewProductCartList struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Weight   int    `json:"weight"`
	Quantity int    `json:"quantity"`
}
