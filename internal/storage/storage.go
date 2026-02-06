package storage

import (
	"errors"
)

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
var ErrRefreshTokenNotFound = errors.New("refresh token not found")
var ErrRefreshTokenExpired = errors.New("refresh token expired")
var ErrInvalidPaginationParams = errors.New("invalid pagination params")
var ErrSubscriptionNotFound = errors.New("subscription not found")
var ErrSubscriptionIDInvalidParams = errors.New("subscription_id params must be an integer")
