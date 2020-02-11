package rate

// Service expose authentication process as a service
type Service interface {
	Authenticate() (*AuthenticationResponse, error)
}

// AuthenticationResponse is the response struct from Spotify
type AuthenticationResponse struct {
	Token      string `json:"access_token"`
	Type       string `json:"token_type"`
	Expiration int    `json:"expires_in"`
	Scope      string `json:"scope"`
}
