package model

type User struct {
	ID       int
	Name     string
	Phone    *string
	Email    string
	Password string
	Bonuses  int
}
