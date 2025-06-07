package model

type Product struct {
	ID            int     `json:"id"`
	CategoryID    int     `json:"category_id"`
	Name          string  `json:"name"`
	Image         string  `json:"image_url"`
	Calorie       float64 `json:"calorie"` // калории
	Weight        int     `json:"weight"`  // вес
	Price         float64 `json:"price"`
	Squirrels     float64 `json:"squirrels"`     // белки
	Fats          float64 `json:"fats"`          // жиры
	Carbohydrates float64 `json:"carbohydrates"` // углеводы
	Visibility    bool
}

type ViewProductList struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	CategoryID int     `json:"category_id"`
	Image      string  `json:"image_url"`
	Weight     int     `json:"weight"`
	Price      float64 `json:"price"`
}

type ViewProductCartList struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Image    string  `json:"image_url"`
	Weight   int     `json:"weight"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
