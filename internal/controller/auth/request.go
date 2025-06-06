package auth

type SignUnRequest struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
