package repository

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicate      = errors.New("record duplication")
	ErrEditConflict   = errors.New("edit conflict")
)
