package customServiceError

import "errors"

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")
var ErrUnknown = errors.New("unknown error. Please, try again")
var ErrCartsNotFound = errors.New("carts not found")
