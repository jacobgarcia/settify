package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SetRequest struct {
	Token string
}

func DecodeAuthRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	keys, ok := r.URL.Query()["token"]

	if !ok || len(keys[0]) < 1 {
		fmt.Println("URL Param token is missing")
		var errResponse IntersectError
		var nestedError NestedError

		nestedError = NestedError{
			Message: "URL param TOKEN is missing",
			Status:  401,
		}

		errResponse = IntersectError{
			Error: nestedError,
		}

		resp, err := json.Marshal(errResponse)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s", resp)
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]

	s := SetRequest{
		Token: key,
	}

	fmt.Println(s)
	return s, nil
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

	if err != nil {
		fmt.Println(err)
		return
	}

	msg := ErrorResponse{
		Message: errResponse.Error.Message,
	}

	w.WriteHeader(errResponse.Error.Status)
	json.NewEncoder(w).Encode(msg)
}
