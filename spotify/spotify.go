package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jacobgarcia/settify/transport"
)

// New instatiates a new API client for Spotify
func New(a, u, i, s string) *Client {
	return &Client{
		authURL: a,
		URL:     u,
		id:      i,
		secret:  s,
	}
}

// Client contains the required params to connect succesfully to Spotify API
type Client struct {
	authURL string
	URL     string
	id      string
	secret  string
}

// Service expose all endpoints as a services.
// This is a microservices architecture
type Service interface {
	Authenticate() (*AuthenticationResponse, error)
	Intersect(token string) (*AuthenticationResponse, error)
}

// AuthenticationResponse is the response struct from Spotify
type AuthenticationResponse struct {
	Token      string `json:"access_token"`
	Type       string `json:"token_type"`
	Expiration int    `json:"expires_in"`
	Scope      string `json:"scope"`
}

var httpClient *http.Client = &http.Client{}
var data url.Values = url.Values{}

// Authenticate is the method for authentication and getting a valid Spotify token
func (c Client) Authenticate() (*AuthenticationResponse, error) {
	uri := fmt.Sprintf("%s/api/token", c.authURL)

	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(data.Encode()))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.id, c.secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		var errResponse transport.ErrorResponse
		err = json.Unmarshal(response, &errResponse)

		if err != nil {
			return nil, err
		}

		errResponse = transport.ErrorResponse{
			Message: errResponse.Message,
			Status:  res.StatusCode,
		}

		resp, err := json.Marshal(errResponse)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s", resp)
	}

	var authResponse AuthenticationResponse
	err = json.Unmarshal(response, &authResponse)

	auth := AuthenticationResponse{
		Token:      authResponse.Token,
		Type:       authResponse.Type,
		Expiration: authResponse.Expiration,
		Scope:      authResponse.Scope,
	}

	return &auth, err
}

// Intersect is the first method will be implementing in Settify. Basically takes two playlists, and generates a new playlist containing the interesection between them.
func (c Client) Intersect(token string) (*AuthenticationResponse, error) {
	uri := fmt.Sprintf("%s/v1/tracks/3n3Ppam7vgaVa1iaRUc9Lp", c.URL)

	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(data.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		var errResponse transport.IntersectError
		err = json.Unmarshal(response, &errResponse)

		if err != nil {
			return nil, err
		}

		errResponse = transport.IntersectError{
			Error: errResponse.Error,
		}

		resp, err := json.Marshal(errResponse)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s", resp)
	}

	var authResponse AuthenticationResponse
	err = json.Unmarshal(response, &authResponse)

	auth := AuthenticationResponse{
		Token:      authResponse.Token,
		Type:       authResponse.Type,
		Expiration: authResponse.Expiration,
		Scope:      authResponse.Scope,
	}

	return &auth, err
}
