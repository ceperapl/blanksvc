package service

import (
	"fmt"

	"github.com/go-kit/log"
)

type Service interface {
	Hello(name string) (string, error)
}

type service struct {
	logger log.Logger
}

func New(logger log.Logger) Service {
	return service{
		logger: logger,
	}
}

func (s service) Hello(name string) (string, error) {
	if name == "" {
		return "", ErrEmplyString
	}

	return fmt.Sprintf("Hello %s", name), nil
}
