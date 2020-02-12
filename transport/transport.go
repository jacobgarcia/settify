package transport

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeAuthRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	// If we would want to parse anything in query/parameteres
	// function also needed for NewServer
	return r, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// ErrorResponse is the standard response message for error handling
type ErrorResponse struct {
	Message string `json:"error"`
	Status  int    `json:"status,omitempty"`
}

// NestedError is the nested response message for error handling
type NestedError struct {
	Message string `json:"message"`
	Status  int    `json:"status,omitempty"`
}

// IntersectError is the standard response message for error handling
type IntersectError struct {
	Error NestedError `json:"error"`
}

// ErrorEncoder returns a REST API response for errors
func ErrorEncoder(c context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var errResponse ErrorResponse

	err = json.Unmarshal([]byte(err.Error()), &errResponse)

	msg := ErrorResponse{
		Message: errResponse.Message,
	}

	w.WriteHeader(errResponse.Status)
	json.NewEncoder(w).Encode(msg)
}

// IntersectErrorEncoder returns a REST API response for errors
func IntersectErrorEncoder(c context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var errResponse IntersectError

	err = json.Unmarshal([]byte(err.Error()), &errResponse)

	msg := ErrorResponse{
		Message: errResponse.Error.Message,
	}

	w.WriteHeader(errResponse.Error.Status)
	json.NewEncoder(w).Encode(msg)
}
