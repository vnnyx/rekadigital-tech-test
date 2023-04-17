package entity

import "errors"

type Customer struct {
	ID        string
	Name      string
	CreatedAt int64
}

var ErrNotFound = errors.New("customer not found")
