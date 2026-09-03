package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/iGuessImaDev/go_lib_response/response"
	"github.com/iGuessImaDev/gocourse_enrollment/internal/enrollment"
)

func NewEnrollmentHTTPServer(ctx context.Context, endpoints enrollment.Endpoints) http.Handler {
	r := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeStoreCourse,
		encodeResponse,
		opts...,
	)).Methods("POST")

	return r
}

func decodeStoreCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var req enrollment.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequestError(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
