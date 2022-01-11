package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"github.com/neidersalgado/go-grpc-api/pkg/entities"
	errorApi "github.com/neidersalgado/go-grpc-api/pkg/errors"
)

type Endpoints struct {
	Authenticate        endpoint.Endpoint
	CreateUserEndpoint  endpoint.Endpoint
	GetUserEndpoint     endpoint.Endpoint
	GetAllUsersEndpoint endpoint.Endpoint
	UpdateUserEndpoint  endpoint.Endpoint
	DeleteUserEndpoint  endpoint.Endpoint
}

func MakeServerEndpoints(s ProxyRepository, logger log.Logger) Endpoints {
	return Endpoints{
		Authenticate:        makeAuthenticateEndpoint(s, logger),
		CreateUserEndpoint:  makeCreateUserEndpoint(s, logger),
		GetUserEndpoint:     makeGetUserEndpoint(s, logger),
		DeleteUserEndpoint:  makeDeleteUserEndpoint(s, logger),
		UpdateUserEndpoint:  makeUpdateUserEndpoint(s, logger),
		GetAllUsersEndpoint: makeGetAlluserEndpoint(s, logger),
	}
}

func makeCreateUserEndpoint(s ProxyRepository, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		reqData, validCast := request.(UserRequest)
		if !validCast {
			return CreateUserResponse{}, errors.New("invalid input data")
		}
		usr := entities.User{
			UserID:                reqData.UserID,
			Email:                 reqData.Email,
			Name:                  reqData.Name,
			Age:                   reqData.Age,
			AdditionalInformation: reqData.AdditionalInformation,
		}
		err = s.Create(usr)
		return CreateUserResponse{Err: err}, nil
	}
}

func makeGetUserEndpoint(s ProxyRepository, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(getUserRequest)
		if !ok {
			return nil, errorApi.NewErrBadRequest(errorApi.ErrInvalidInputType)
		}
		user, err := s.Get(ctx, req.Email)
		if err != nil {
			return nil, err
		}

		return UserResponse{
			Email:                 user.Email,
			AdditionalInformation: user.AdditionalInformation,
			Age:                   user.Age,
			Name:                  user.Name,
		}, nil
	}
}

func makeGetAlluserEndpoint(s ProxyRepository, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, err := s.List(ctx)
		if err != nil {
			return nil, err
		}

		return getAllResponse{Users: []UserResponse{}}, nil
	}
}

func makeDeleteUserEndpoint(s ProxyRepository, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(getUserRequest)
		if !ok {
			return nil, errorApi.NewErrBadRequest(errorApi.ErrInvalidInputType)
		}
		err := s.Delete(ctx, req.Email)
		if err != nil {
			return DeleteUserResponse{Err: err}, err
		}

		return DeleteUserResponse{Msg: "user deleted OK"}, nil
	}
}

func makeUpdateUserEndpoint(s ProxyRepository, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		reqData, validCast := request.(UserRequest)
		if !validCast {
			return nil, errors.New("invalid input data")
		}
		usr := entities.User{
			UserID:                reqData.UserID,
			Email:                 reqData.Email,
			Name:                  reqData.Name,
			Age:                   reqData.Age,
			AdditionalInformation: reqData.AdditionalInformation,
		}
		err = s.Update(ctx, usr)

		if err != nil {
			fmt.Printf("Error: %+v", err.Error())
			return UpdateUserResponse{Err: err}, err
		}

		return UpdateUserResponse{}, nil
	}
}

func makeAuthenticateEndpoint(s ProxyRepository, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(authRequest)
		if !ok {
			return nil, errorApi.NewErrBadRequest(errorApi.ErrInvalidInputType)
		}
		auth, err := s.Authenticate(req.Email, req.Hash)
		if err != nil || !auth {
			return AuthResponse{
				Err: err,
				Msg: "User couldn't be Authenticated",
			}, err
		}

		return AuthResponse{
			Err: nil,
			Msg: "User Well Authenticated",
		}, nil
	}
}
