package profilesvc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(ctx context.Context, s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	r.Methods("POST").Path("/profile").Handler(httptransport.NewServer(
		ctx,
		MakePostProfileEndpoint(s),
		decodePostProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/profile/{id}").Handler(httptransport.NewServer(
		ctx,
		MakeGetProfileEndpoint(s),
		decodeGetProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/profile/{id}").Handler(httptransport.NewServer(
		ctx,
		MakeDeleteProfileEndpoint(s),
		decodeDeleteProfileRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodePostProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)

	var request postProfileRequest
	err := json.NewDecoder(body).Decode(&request.Profile)
	if err != nil {
		return nil, err
	}
	fmt.Println("body", buf.String())
	return request, nil
}

func decodeGetProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getProfileRequest{ID: id}, nil
}

func decodeDeleteProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteProfileRequest{ID: id}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}