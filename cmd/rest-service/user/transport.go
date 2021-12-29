package user

import (
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	ErrRouting = errors.New("inconsistent mapping route")
)

// MakeHTTPHandler set up services in handlers
func MakeHTTPHandler(s ProxyRepository, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodPost).Path(PostUser).Handler(httptransport.NewServer(
		e.CreateUserEndpoint,
		decodeCreateRequest,
		encodeCreateResponse,
		options...,
	))

	r.Methods(http.MethodGet).Path(GetUser).Handler(httptransport.NewServer(
		e.GetUserEndpoint,
		decoderEmailRequest,
		encodeCreateResponse,
		options...,
	))
	r.Methods(http.MethodDelete).Path(DeleteUser).Handler(httptransport.NewServer(
		e.DeleteUserEndpoint,
		decoderEmailRequest,
		encodeDeleteResponse,
		options...,
	))
	r.Methods(http.MethodPut).Path(UpdateUser).Handler(httptransport.NewServer(
		e.UpdateUserEndpoint,
		decodeUpdateRequest,
		encodeDeleteResponse,
		options...,
	))

	return r
}
