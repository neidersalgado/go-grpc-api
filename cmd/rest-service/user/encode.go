package user

import (
	"context"
	"encoding/json"
	"net/http"

	apiErrors "github.com/neidersalgado/go-grpc-api/pkg/errors"
)

func encodeCreateResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {

	if err == nil {

		panic("encodeError with nil error")

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})

}

func codeFrom(err error) int {
	switch err.(type) {
	case apiErrors.ErrNotFound:
		return http.StatusNotFound
	case apiErrors.ErrBadRequest:
		return http.StatusBadRequest
	case apiErrors.ErrForbidden:
		return http.StatusForbidden
	case apiErrors.ErrInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func encodeDeleteResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, validCast := response.(DeleteUserResponse)

	if !validCast {
		w.WriteHeader(http.StatusBadRequest)
		resp.Err = apiErrors.NewErrBadRequest("delete response could not be read")
		return json.NewEncoder(w).Encode(resp)
	}

	if resp.Err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(w).Encode(resp)
	}
	w.WriteHeader(http.StatusOK)
	resp.Msg = "Deleted Ok"

	return json.NewEncoder(w).Encode(resp)
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}
