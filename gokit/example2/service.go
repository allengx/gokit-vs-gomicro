package main

import (
	"context"
	"errors"
	"strings"
)

// option interface about string
type StringService interface {
	Uppercase(context.Context, string) (string, error)
}

// option struct
type stringService struct{}

// func realize
func (str stringService) Uppercase(ctx context.Context, s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

var ErrEmpty = errors.New("Empty string")
