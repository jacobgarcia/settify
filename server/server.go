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

var logger log.Logger
var service spotify.Service

// CreateRouter is in charge to define all routes
func CreateRouter(spotifyService spotify.Service, serverLogger log.Logger) http.Handler {
	// Set the logger we will be using in the server
	logger = serverLogger
	service = spotifyService
	r := mux.NewRouter()

	playlistsHandler := getHandler(playlistsEndpoint())
	intersectHandler := getHandler(operationEndpoint("intersection"))
	unionHandler := getHandler(operationEndpoint("union"))
	profileHandler := getHandler(profileEndpoint())
	complementHandler := getHandler(operationEndpoint("complement"))
	userPlaylistsHandler := getHandler(usersEndpoint())

	// Basic Spotify calls
	r.Handle("/me", profileHandler).Methods("GET")
	r.Handle("/playlists", playlistsHandler).Methods("GET")
	r.Handle("/user/playlists", userPlaylistsHandler).Methods("GET")
	// Set operations
	r.Handle("/intersection", intersectHandler).Methods("GET")
	r.Handle("/union", unionHandler).Methods("GET")
	r.Handle("/complement", complementHandler).Methods("GET")
	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(r)
}

func playlistsEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transport.AuthRequest)
		auth, err := service.Playlists(req.Token, req.Offset)
		if err != nil {
			return nil, err
		}
		return auth, nil
	}
}

func profileEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transport.AuthRequest)
		auth, err := service.Profile(req.Token)
		if err != nil {
			return nil, err
		}
		return auth, nil
	}
}

func usersEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transport.AuthRequest)
		auth, err := service.UserPlaylists(req.Token, req.Offset, req.Username)
		if err != nil {
			return nil, err
		}
		return auth, nil
	}
}

func operationEndpoint(operation string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transport.AuthRequest)
		auth, err := &spotify.NewPlaylistResponse{}, nil

		switch operation {
		case "intersection":
			auth, err = service.Intersect(req.Token, req.FirstPlaylist, req.SecondPlaylist, req.Name)
		case "union":
			auth, err = service.Union(req.Token, req.FirstPlaylist, req.SecondPlaylist, req.Name)
		case "complement":
			auth, err = service.Complement(req.Token, req.FirstPlaylist, req.SecondPlaylist, req.Name)
		}
		if err != nil {
			return auth, err
		}
		return auth, err
	}
}

func getHandler(endpoint endpoint.Endpoint) *kithttp.Server {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(transport.IntersectErrorEncoder),
	}

	return kithttp.NewServer(
		endpoint,
		transport.DecodeAuthRequest,
		transport.EncodeResponse,
		opts...)
}
