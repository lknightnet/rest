package model

type User struct {
	ID       int
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string
	Bonuses  int     `json:"bonuses"`
	Address  *string `json:"address"`
}

type ViewUser struct {
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	Bonuses int     `json:"bonuses"`
	Address *string `json:"address"`
}
