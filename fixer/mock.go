package fixer

import (
	"github.com/adrianforsius/go-service/rate"
)

type Mock struct{}

func (m Mock) LatestRate() (*rate.Provider, error) {
	return &rate.Provider{}, nil
}
