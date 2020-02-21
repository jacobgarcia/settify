package server

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/jacobgarcia/settify/spotify"
	"github.com/jacobgarcia/settify/transport"
)

// CreateRouter is in charge to define all routes
func CreateRouter(s spotify.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(transport.IntersectErrorEncoder),
	}

	playlistsHandler := kithttp.NewServer(
		playlistsEndpoint(s, logger),
		transport.DecodeAuthRequest,
		transport.EncodeResponse,
		opts...)

	intersectHandler := kithttp.NewServer(
		intersectEndpoint(s, logger, "intersection"),
		transport.DecodeAuthRequest,
		transport.EncodeResponse,
		opts...)

	unionHandler := kithttp.NewServer(
		intersectEndpoint(s, logger, "union"),
		transport.DecodeAuthRequest,
		transport.EncodeResponse,
		opts...)

	r.Handle("/playlists", playlistsHandler).Methods("GET")
	r.Handle("/intersection", intersectHandler).Methods("GET")
	r.Handle("/union", unionHandler).Methods("GET")
	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(r)
}

func playlistsEndpoint(service spotify.Service, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transport.AuthRequest)
		auth, err := service.Playlists(req.Token, req.Offset)
		if err != nil {
			return auth, err
		}
		return auth, err
	}
}

func intersectEndpoint(service spotify.Service, logger log.Logger, operation string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transport.AuthRequest)
		auth, err := &spotify.NewPlaylistResponse{}, nil

		switch operation {
		case "intersection":
			auth, err = service.Intersect(req.Token, req.FirstPlaylist, req.SecondPlaylist)
		case "union":
			auth, err = service.Union(req.Token, req.FirstPlaylist, req.SecondPlaylist)
		}
		if err != nil {
			return auth, err
		}
		return auth, err
	}
}
