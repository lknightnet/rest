package model

type Product struct {
	ID            int     `json:"id"`
	CategoryID    int     `json:"category_id"`
	Name          string  `json:"name"`
	Image         string  `json:"image_url"`
	Calorie       float64 `yaml:"calorie"` // калории
	Weight        int     `yaml:"weight"`  // вес
	Price         float64 `json:"price"`
	Squirrels     float64 `yaml:"squirrels"`     // белки
	Fats          float64 `yaml:"fats"`          // жиры
	Carbohydrates float64 `yaml:"carbohydrates"` // углеводы
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
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image_url"`
	Weight   int    `json:"weight"`
	Quantity int    `json:"quantity"`
}
