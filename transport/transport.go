package transport

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeGetRateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	// If we would want to parse anything in query/parameteres
	// function also needed for NewServer
	return r, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
