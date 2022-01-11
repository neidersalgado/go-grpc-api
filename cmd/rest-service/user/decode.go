package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

var errBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")

func decodeCreateRequest(logger log.Logger) func(_ context.Context, r *http.Request) (request interface{}, err error) {
	return func(_ context.Context, r *http.Request) (request interface{}, err error) {
		logger.Log("decode", "decode request")
		var req UserRequest
		if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
			logger.Log("decode", "Decode Failed")
			return nil, e
		}
		logger.Log("decode", "decode Ok")
		return req, nil
	}
}

func decoderEmailRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		return nil, errBadRouting
	}
	return getUserRequest{Email: email}, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req UserRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		return nil, errBadRouting
	}
	req.Email = email

	return req, nil
}
