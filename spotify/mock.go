package spotify

import (
	"github.com/jacobgarcia/settify/rate"
)

type Mock struct{}

func (m Mock) Authenticate() (*rate.AuthenticationResponse, error) {
	return &rate.AuthenticationResponse{}, nil
}
