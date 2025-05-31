package model

type Category struct {
	ID         int
	Name       string
	ImageUrl   string
	Sort       int
	Visibility bool
}

type ViewCategoryList struct {
	ID       int
	Name     string
	ImageUrl string
	Sort     int
}

type ViewCategoryWithProductList struct {
	ID          int
	Name        string
	Sort        int
	ImageUrl    string
	ProductList []ViewProductList
}
