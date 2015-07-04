package models

import (
	"fmt"
)

const (
	ErrInvalidId ModelError = iota + 1
	ErrNotFound
	ErrDuplicateKey
)

type ModelError int

func (e ModelError) Error() string {
	switch e {
	case ErrInvalidId:
		return "invalid object id"
	case ErrDuplicateKey:
		return "duplicate key"
	case ErrNotFound:
		return "not found"
	default:
		return fmt.Sprintf("undefined model error, number: %d", int(e))
	}
}
