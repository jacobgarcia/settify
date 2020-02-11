package spotify

type Mock struct{}

func (m Mock) Authenticate() (*AuthenticationResponse, error) {
	return &AuthenticationResponse{}, nil
}
