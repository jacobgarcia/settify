package fixer

import (
	"github.com/jacobgarcia/settify/rate"
)

type Mock struct{}

func (m Mock) LatestRate() (*rate.Provider, error) {
	return &rate.Provider{}, nil
}
