package user

type ChangeRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Address  string `json:"address"`
}
