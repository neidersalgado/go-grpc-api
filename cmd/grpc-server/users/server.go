package grpc

import (
	"context"
	"errors"

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
		getUser:      grpc.NewServer(endpoints.GetUserEndpoint, decodeGetUserRequest, encodeGetUserResponse, options...),
		authenticate: grpc.NewServer(endpoints.AuthenticateUserEndpoint, nil, nil, options...),
		createUser:   grpc.NewServer(endpoints.CreateUserEndpoint, decodeCreateUserRequest, encodeCreateUserResponse, options...),
		getAll:       grpc.NewServer(endpoints.getAllEndpoint, nil, nil, options...),
		update:       grpc.NewServer(endpoints.UpdateUserEndpoint, nil, nil, options...),
		delete:       grpc.NewServer(endpoints.DeleteUserEndpoint, nil, nil, options...),
	}

	return server
}

func (srv *grpcUserServer) Authenticate(ctx context.Context, auth *pb.AuthRequest) (*pb.Response, error) {
	_, grpcResponse, err := srv.authenticate.ServeGRPC(ctx, auth)

	return grpcResponse.(*pb.Response), err
}

func (srv *grpcUserServer) Create(ctx context.Context, user *pb.UserRequest) (*pb.Response, error) {
	_, grpcResponse, err := srv.createUser.ServeGRPC(ctx, user)

	response, ok := grpcResponse.(*pb.Response)
	if !ok {
		return nil, errors.New("invalid input data create")
	}
	return response, err
}

func (srv *grpcUserServer) Get(ctx context.Context, uid *pb.UserIDRequest) (*pb.UserResponse, error) {
	_, grpcResponse, err := srv.getUser.ServeGRPC(ctx, uid)

	return grpcResponse.(*pb.UserResponse), err
}

func (srv *grpcUserServer) Update(ctx context.Context, user *pb.UserRequest) (*pb.Response, error) {
	_, grpcResponse, err := srv.update.ServeGRPC(ctx, user)

	return grpcResponse.(*pb.Response), err
}

func (srv *grpcUserServer) Delete(ctx context.Context, uid *pb.UserIDRequest) (*pb.Response, error) {
	_, grpcResponse, err := srv.delete.ServeGRPC(ctx, uid)

	return grpcResponse.(*pb.Response), err
}

func (srv *grpcUserServer) GetAll(ctx context.Context, void *pb.Void) (*pb.UserColletionResponse, error) {
	_, grpcResponse, err := srv.getAll.ServeGRPC(ctx, void)

	return grpcResponse.(*pb.UserColletionResponse), err
}
