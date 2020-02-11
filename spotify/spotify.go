package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jacobgarcia/settify/rate"
	"github.com/jacobgarcia/settify/transport"
)

func New(a, u, i, s string) *client {
	return &client{
		authURL: a,
		URL:     u,
		id:      i,
		secret:  s,
	}
}

type client struct {
	authURL string
	URL     string
	id      string
	secret  string
}

func (c client) Authenticate() (*rate.AuthenticationResponse, error) {
	client := &http.Client{}
	uri := fmt.Sprintf("%s/api/token", c.authURL)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(data.Encode()))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.id, c.secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)

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

	var authResponse rate.AuthenticationResponse
	err = json.Unmarshal(response, &authResponse)

	auth := rate.AuthenticationResponse{
		Token:      authResponse.Token,
		Type:       authResponse.Type,
		Expiration: authResponse.Expiration,
		Scope:      authResponse.Scope,
	}

	return &auth, err
}
