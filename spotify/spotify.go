package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/jacobgarcia/settify/rate"
	"github.com/jacobgarcia/settify/transport"
)

func New(u, i, s string) *client {
	return &client{
		URL:    u,
		id:     i,
		secret: s,
	}
}

type client struct {
	URL    string
	id     string
	secret string
}

type ListResponse struct {
	Timestamp int64     `json:"timestamp"`
	Rates     rate.Rate `json:"rates"`
}

func (c client) LatestRate() (*rate.Provider, error) {
	client := &http.Client{}
	uri := fmt.Sprintf("%s/api/token", c.URL)

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

	var list ListResponse
	err = json.Unmarshal(response, &list)

	rate := rate.Provider{
		Rate:    list.Rates.MXN,
		Updated: time.Unix(list.Timestamp, 0).Format(time.RFC3339),
		Name:    "fixer",
	}

	return &rate, err
}
