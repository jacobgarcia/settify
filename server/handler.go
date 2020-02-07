package server

import (
	"encoding/json"
	"net/http"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/adrianforsius/go-service/rate"
)

func MakeHTTPHandler(s rate.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	key := []byte("supersecret")
	keys := func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerBefore(stdjwt.ToHTTPContext()),
	}

	getRateHandler := kithttp.NewServer(
		stdjwt.NewParser(keys, jwt.SigningMethodHS256, &customClaims{})(ratesEndpoint(s, logger)),
		DecodeGetRateRequest,
		EncodeResponse,
		opts...,
	)

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(authErrorEncoder),
		httptransport.ServerErrorLogger(logger),
	}
	authHandler := httptransport.NewServer(
		makeAuthEndpoint(auth),
		decodeAuthRequest,
		encodeResponse,
		opts...,
	)
	http.Handle("/auth", methodControl("POST", authHandler))

	r.Handle("/rates", getRateHandler).Methods("GET")
	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	return r
}

func ratesEndpoint(r rate.Service, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		logger.Log("rates endpoint hit")
		provider, err := r.LatestRate()
		if err != nil {
			return provider, err
		}
		return ProviderResponse{
			Providers: []rate.Provider{
				*provider,
			},
		}, err
	}
}

func authEndpoint(svc AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(authRequest)
		token, err := svc.Auth(req.ClientID, req.ClientSecret)
		if err != nil {
			return nil, err
		}
		return authResponse{token, ""}, nil
	}
}

type ProviderResponse struct {
	Providers []rate.Provider
}
