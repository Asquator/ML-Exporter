package storage

import "errors"

var (
	ErrModelNotFound = errors.New("model not found")
	ErrModelExists   = errors.New("model already exists")
)
