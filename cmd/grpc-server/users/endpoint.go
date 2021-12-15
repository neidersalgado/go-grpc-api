package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"github.com/neidersalgado/go-camp-grpc/cmd/grpc-server/pb"
	"github.com/neidersalgado/go-camp-grpc/pkg/users"
	domain "github.com/neidersalgado/go-camp-grpc/pkg/users"
)

const (
	invalidData = "invalid request data %s"
	invalidAuth = "no user found"
	notCreated  = "User couldn't be created \n %v"
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
			return nil, errors.New(fmt.Sprintf(invalidData, "Authenticate"))
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
		requestData, validCast := request.(createUserRequest)

		if !validCast {
			return nil, errors.New(fmt.Sprintf(invalidData, ": Create"))
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
			return nil, errors.New(fmt.Sprintf(notCreated, err))
		}

		return createUserResponse{Id: usr.Email, Error: err}, nil
	}
}

func MakeGetUserEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestData, validCast := request.(getUserRequest)

		if !validCast {
			return nil, errors.New(fmt.Sprintf(invalidData, " get"))
		}

		usr, err := s.GetByEmail(ctx, requestData.Email)

		return getUserResponse{
			UserResponse: UserResponse{
				UserId:                usr.UserId,
				PwdHash:               usr.PwdHash,
				Email:                 usr.Email,
				Name:                  usr.Name,
				Age:                   usr.Age,
				AdditionalInformation: usr.AdditionalInformation ,
			},
		}, nil
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
