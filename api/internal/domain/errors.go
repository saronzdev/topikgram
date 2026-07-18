package domain

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrIncorrectPassword = errors.New("incorrect password")
)