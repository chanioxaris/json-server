package storage

import (
	"errors"
)

var (
	// ErrResourceNotFound returns an error when a requested resource not found in storage.
	ErrResourceNotFound = errors.New("resource not found")
	// ErrResourceAlreadyExists returns an error when a resource already exists in storage.
	ErrResourceAlreadyExists = errors.New("resource already exists")
	// ErrBadRequest returns an error when an unexpected request been processed.
	ErrBadRequest = errors.New("bad request")
	// ErrInternalServerError returns an error when an unexpected error occurs.
	ErrInternalServerError = errors.New("internal Server Error")
)

// Resource represents the structure of a singe resource in storage.
type Resource map[string]interface{}

// Service interface to handle storage operations.
type Service interface {
	Find() ([]Resource, error)
	FindById(string) (Resource, error)
	Create(Resource) (Resource, error)
	Replace(string, Resource) (Resource, error)
	Delete(string) error
}
