package service

import "errors"

var (
	ErrFailedValidation = errors.New("validation failed")
	ErrWrongCredentials = errors.New("wrong user credentials")
	ErrDuplicate        = errors.New("record duplication")
)
