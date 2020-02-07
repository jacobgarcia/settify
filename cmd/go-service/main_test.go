package main

import (
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/valyala/fasthttp"

	"github.com/adrianforsius/go-service/fixer"
	"github.com/adrianforsius/go-service/server"
)

func TestService(t *testing.T) {
	ts := Setup(t)
	defer TearDown(ts)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(ts.URL + "/healthcheck")

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		t.Errorf("Unexpected error calling the service: %s", err)
	}

	if resp.StatusCode() != 200 {
		t.Errorf("Expected %d, Got %s", 200, resp.StatusCode())
	}

	req.SetRequestURI(ts.URL + "/rates")
	err = client.Do(req, resp)
	if err != nil {
		t.Errorf("Unexpected error calling the service: %s", err)
	}

	// Rates are dynamic we need a mock here to check the response body
	if resp.StatusCode() != 200 {
		t.Errorf("Expected %d, Got %s", 200, resp.StatusCode())
	}

	// TODO: this is getting repetative use a helper function
	req.SetRequestURI(ts.URL + "/not-exist")
	err = client.Do(req, resp)
	if err != nil {
		t.Errorf("Unexpected error calling the service: %s", err)
	}

	if resp.StatusCode() != 404 {
		t.Errorf("Expected %d, Got %s", 404, resp.StatusCode())
	}
}

func Setup(t *testing.T) *httptest.Server {
	t.Helper()
	handler := server.MakeHTTPHandler(fixer.Mock{}, log.NewNopLogger())
	return httptest.NewServer(handler)
}

func TearDown(ts *httptest.Server) {
	defer ts.Close()
}
