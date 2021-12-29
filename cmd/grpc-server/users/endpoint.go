package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	domain "github.com/neidersalgado/go-grpc-api/pkg/users"
)

const (
	invalidData = "invalid request data %v"
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
		requestData, validCast := request.(domain.Auth)
		if !validCast {
			return Response{Code: 400}, errors.New(fmt.Sprintf(invalidData, "Authenticate"))
		}

		err = s.Authenticate(ctx, requestData)
		if err != nil {
			return response, err
		}

		return Response{Code: 200}, nil
	}
}

func MakeCreateUserEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestData, validCast := request.(createUserRequest)

		if !validCast {
			return createUserResponse{}, errors.New(fmt.Sprintf(invalidData, ": Create"))
		}
		usr := domain.User{
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
		requestData, validCast := request.(userIdRequest)

		if !validCast {
			return getAllUsersResponse{}, errors.New(fmt.Sprintf(invalidData, " get"))
		}

		usr, err := s.GetByEmail(ctx, requestData.Email)

		return getUserResponse{
			UserResponse: UserResponse{
				UserId:                usr.UserId,
				PwdHash:               usr.PwdHash,
				Email:                 usr.Email,
				Name:                  usr.Name,
				Age:                   usr.Age,
				AdditionalInformation: usr.AdditionalInformation,
			},
		}, nil
	}
}

func MakeGetAllEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		usersData, err := s.GetAll(ctx)

		allUsersResponse := getAllUsersResponse{Users: []UserResponse{}}

		for _, usr := range usersData {
			UserResponse := UserResponse{
				UserId:                usr.UserId,
				PwdHash:               usr.PwdHash,
				Email:                 usr.Email,
				Name:                  usr.Name,
				Age:                   usr.Age,
				AdditionalInformation: usr.AdditionalInformation,
			}
			allUsersResponse.Users = append(allUsersResponse.Users, UserResponse)
		}

		return allUsersResponse, nil
	}
}

func MakeUpdateUserEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestData, validCast := request.(updateUserRequest)
		fmt.Printf("cast: %+v", validCast)
		if !validCast {
			return Response{Code: 400}, errors.New(fmt.Sprintf(invalidData, "Update"))
		}

		usr := domain.User{
			UserId:                requestData.UserId,
			Email:                 requestData.Email,
			PwdHash:               requestData.PwdHash,
			Name:                  requestData.Name,
			Age:                   requestData.Age,
			AdditionalInformation: requestData.AdditionalInformation,
		}
		err = s.Update(ctx, usr)
		if err != nil {
			return Response{Code: 400}, errors.New(fmt.Sprintf(invalidData, err))
		}
		return Response{Code: 200}, nil
	}
}

func MakeDeleteUserEndpoint(s domain.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestData, validCast := request.(userIdRequest)

		if !validCast {
			return Response{Code: 400}, errors.New(invalidData)
		}
		err = s.Delete(ctx, requestData.Email)
		if err != nil {
			return Response{Code: 400}, err
		}
		return Response{Code: 200}, nil
	}
}
