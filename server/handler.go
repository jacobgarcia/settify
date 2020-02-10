package server

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/jacobgarcia/settify/rate"
)

func MakeHTTPHandler(s rate.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()

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
