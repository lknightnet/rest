package customRepositoryError

import "errors"

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")
var ErrUserNotFound = errors.New("user not found")
var ErrProductNotFound = errors.New("product not found")
var ErrCartNotFound = errors.New("cart not found")
var ErrCartsNotFound = errors.New("carts not found")
