package storage

import "errors"

var (
	ErrNotFound  = errors.New("url not found")
	ErrUrlExists = errors.New("url already exists")
)
