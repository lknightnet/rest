package user

type ErrorResponse struct {
	Error string `json:"error"`
}

type OKResponse struct {
	Status bool `json:"status"`
}
