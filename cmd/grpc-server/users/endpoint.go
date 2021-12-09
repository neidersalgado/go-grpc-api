package grpc

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/neidersalgado/go-camp-grpc/cmd/grpc-server/pb"
	"github.com/neidersalgado/go-camp-grpc/pkg/users"
	domain "github.com/neidersalgado/go-camp-grpc/pkg/users"
)

const (
	invalidData = "invalid request data"
	invalidAuth = "no user found"
	notCreated  = "User couldn't be created"
)

type grpcUserEndpoints struct {
	AuthenticateUserEndpoint endpoint.Endpoint
	CreateUserEndpoint       endpoint.Endpoint
	GetUserEndpoint          endpoint.Endpoint
	getAllEndpoint           endpoint.Endpoint
	UpdateUserEndpoint       endpoint.Endpoint
	DeleteUserEndpoint       endpoint.Endpoint
}

func NewGrpcUserServerEndpoints(s domain.UserService) *grpcUserEndpoints {
	return &grpcUserEndpoints{
		AuthenticateUserEndpoint: MakeAuthenticateUserEndpoint(s),
		CreateUserEndpoint:       MakeCreateUserEndpoint(s),
		GetUserEndpoint:          MakeGetUserEndpoint(s),
		getAllEndpoint:           MakeGetAllEndpoint(s),
		UpdateUserEndpoint:       MakeUpdateUserEndpoint(s),
		DeleteUserEndpoint:       MakeDeleteUserEndpoint(s),
	}
}

func MakeAuthenticateUserEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestData, validCast := request.(pb.UserIDRequest)

		if !validCast {
			return nil, errors.New(invalidData)
		}

		usr, err := s.GetByEmail(ctx, requestData.Email)

		if usr.UserId >= 0 {
			return nil, errors.New(invalidAuth)
		}

		return usr, nil
	}
}

func MakeCreateUserEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestData, validCast := request.(UserRequest)

		if !validCast {
			return nil, errors.New(invalidData)
		}
		usr := users.User{
			UserId:                requestData.UserId,
			Email:                 requestData.Email,
			PwdHash:               requestData.PwdHash,
			Name:                  requestData.Name,
			Age:                   requestData.Age,
			AdditionalInformation: requestData.AdditionalInformation,
		}
		err = s.Create(ctx, usr)

		if err != nil {
			return nil, errors.New(notCreated)
		}

		return createUserResponse{Id: usr.Email, Error: err}, nil
	}
}

func MakeGetUserEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestData, validCast := request.(pb.UserIDRequest)

		if !validCast {
			return nil, errors.New(invalidData)
		}

		response, err = s.GetByEmail(ctx, requestData.Email)

		return
	}
}

func MakeGetAllEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_, validCast := request.(pb.Void)

		if !validCast {
			return nil, errors.New(invalidData)
		}

		response, err = s.GetAll(ctx)

		return
	}
}

func MakeUpdateUserEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestData, validCast := request.(pb.UserRequest)

		if !validCast {
			return nil, errors.New(invalidData)
		}

		usr := users.User{
			UserId:                requestData.UserId,
			Email:                 requestData.Email,
			PwdHash:               requestData.PwdHash,
			Name:                  requestData.Name,
			Age:                   requestData.Age,
			AdditionalInformation: requestData.AdditionalInformation,
		}
		err = s.Update(ctx, usr)
		return
	}
}

func MakeDeleteUserEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestData, validCast := request.(pb.UserIDRequest)

		if !validCast {
			return nil, errors.New(invalidData)
		}

		err = s.Delete(ctx, requestData.Email)
		if err != nil {

		}
		return
	}
}
