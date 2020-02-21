package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AuthRequest is an authenticated request containing token
type AuthRequest struct {
	Token          string
	FirstPlaylist  string
	SecondPlaylist string
	Offset         string
}

// DecodeAuthRequest serves as a middleware function to intercept requests in order to get the Authorization Bearer Token
func DecodeAuthRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	token := req.Header.Get("Authorization")

	offset := req.URL.Query().Get("offset")
	firstPlaylist := req.URL.Query().Get("firstPlaylist")
	secondPlaylist := req.URL.Query().Get("secondPlaylist")

	if token == "" {
		fmt.Println("Token is missing")
		var errResponse IntersectError
		var nestedError NestedError

		nestedError = NestedError{
			Message: "Bearer TOKEN is missing",
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

	s := AuthRequest{
		Token:          token,
		FirstPlaylist:  firstPlaylist,
		SecondPlaylist: secondPlaylist,
		Offset:         offset,
	}

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
