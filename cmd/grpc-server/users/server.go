package grpc

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"

	"github.com/neidersalgado/go-camp-grpc/cmd/grpc-server/pb"
)

type grpcUserServer struct {
	pb.UsersServer
	getUser      grpc.Handler
	createUser   grpc.Handler
	getAll       grpc.Handler
	update       grpc.Handler
	delete       grpc.Handler
	authenticate grpc.Handler
}

func NewGrpcUserServer(endpoints grpcUserEndpoints, logger log.Logger) pb.UsersServer {
	options := []grpc.ServerOption{
		grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}
	server := &grpcUserServer{
		getUser:      grpc.NewServer(endpoints.GetUserEndpoint, nil, nil, options...),
		authenticate: grpc.NewServer(endpoints.AuthenticateUserEndpoint, nil, nil, options...),
		createUser:   grpc.NewServer(endpoints.CreateUserEndpoint, nil, nil, options...),
		getAll:       grpc.NewServer(endpoints.getAllEndpoint, nil, nil, options...),
		update:       grpc.NewServer(endpoints.UpdateUserEndpoint, nil, nil, options...),
		delete:       grpc.NewServer(endpoints.DeleteUserEndpoint, nil, nil, options...),
	}

	return server
}
