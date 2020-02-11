package server

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/jacobgarcia/settify/rate"
	"github.com/jacobgarcia/settify/transport"
)

// CreateRouter is in charge to define all routes
func CreateRouter(s rate.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(transport.ErrorEncoder),
	}

	getRateHandler := kithttp.NewServer(
		ratesEndpoint(s, logger),
		transport.DecodeGetRateRequest,
		transport.EncodeResponse,
		opts...,
	)

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

type ProviderResponse struct {
	Providers []rate.Provider
}
