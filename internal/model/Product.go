package model

type Product struct {
	ID            int
	CategoryID    int
	Name          string
	Image         string
	Calorie       float64 // калории
	Weight        int     // вес
	Price         float64
	Squirrels     float64 // белки
	Fats          float64 // жиры
	Carbohydrates float64 // углеводы
	Visibility    bool
}

type ViewProductList struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	CategoryID int     `json:"category_id"`
	Image      string  `json:"image"`
	Weight     int     `json:"weight"`
	Price      float64 `json:"price"`
}

type ViewProductCartList struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Weight   int    `json:"weight"`
	Quantity int    `json:"quantity"`
}
