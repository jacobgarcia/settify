package fixer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jacobgarcia/settify/rate"
)

func New(u, t string) *client {
	return &client{
		URL:   u,
		token: t,
	}
}

type client struct {
	URL   string
	token string
}

type ListResponse struct {
	Timestamp int64     `json:"timestamp"`
	Rates     rate.Rate `json:"rates"`
}

func (c client) LatestRate() (*rate.Provider, error) {
	url := fmt.Sprintf("%s/latest?access_key=%s", c.URL, c.token)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 (%d), error: %s\n", resp.StatusCode, response)
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
