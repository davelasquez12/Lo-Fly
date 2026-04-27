package tracker

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("tracked flight not found")
	ErrInactive       = errors.New("tracked flight is not active")
)
