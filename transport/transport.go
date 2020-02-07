package transport

func DecodeGetRateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	// If we would want to parse anything in query/parameteres
	// function also needed for NewServer
	return r, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func DecodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request authRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeAuthError(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusUnauthorized
	msg := err.Error()

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(authResponse{Token: "", Err: msg})
}
