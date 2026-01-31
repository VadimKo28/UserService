package storage

import "errors"

const (
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusUnprocessableEntity = 422
	StatusInternalServerError = 500
)

var ErrUserNotFount = errors.New("record not found")
var ErrUserAlreadyExists = errors.New("user with this email already exists")
var ErrInternalServerError = errors.New("internal server error")
var ErrUserIDInvalidParams = errors.New("id params must be an integer")
var ErrUserInvalidCredentials = errors.New("invalid email or password")
var ErrUnauthorized = errors.New("authorization failed")
var ErrForbidden = errors.New("forbidden")
