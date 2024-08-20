package errors

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

const (
	ErrUnauthorized          = "unauthorized"
	ErrAuthHeaderMissing     = "auth header is missing"
	ErrAuthHeaderInvalid     = "invalid authorization header format"
	ErrTokenSignatureInvalid = "token signature is invalid"
	ErrTokenInvalid          = "invalid authorization token"

	ErrInvalidRequest = "invalid request"
	ErrInternalServer = "internal server error"

	ErrInvalidCredentials = "invalid credentials"
)
