package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/neidersalgado/go-camp-grpc/cmd/user/pb"
)

type grpcUserEndpoints struct {
	AuthenticateUserEndpoint endpoint.Endpoint
	CreateUserEndpoint       endpoint.Endpoint
	GetUserEndpoint          endpoint.Endpoint
	UpdateUserEndpoint       endpoint.Endpoint
	DeleteUserEndpoint       endpoint.Endpoint
}

func NewGrpcUserServer(s pb.UsersServer) *grpcUserEndpoints {
	return &grpcUserEndpoints{
		AuthenticateUserEndpoint: MakeAuthenticateUserEndpoint(s),
		CreateUserEndpoint:       MakeCreateUserEndpoint(s),
		GetUserEndpoint:          MakeGetUserEndpoint(s),
		UpdateUserEndpoint:       MakeUpdateUserEndpoint(s),
		DeleteUserEndpoint:       MakeDeleteUserEndpoint(s),
	}
}

func MakeAuthenticateUserEndpoint(s pb.UsersServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

	}
}

func MakeCreateUserEndpoint(s pb.UsersServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

	}
}

func MakeGetUserEndpoint(s pb.UsersServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

	}
}

func MakeUpdateUserEndpoint(s pb.UsersServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

	}
}

func MakeDeleteUserEndpoint(s pb.UsersServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

	}
}
