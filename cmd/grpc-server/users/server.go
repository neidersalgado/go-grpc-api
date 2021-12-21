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
		getUser:      grpc.NewServer(endpoints.GetUserEndpoint, decodeUserIdRequest, encodeGetUserResponse, options...),
		authenticate: grpc.NewServer(endpoints.AuthenticateUserEndpoint, decodeAuthUserRequest, encodeResponse, options...),
		createUser:   grpc.NewServer(endpoints.CreateUserEndpoint, decodeCreateUserRequest, encodeCreateUserResponse, options...),
		getAll:       grpc.NewServer(endpoints.getAllEndpoint, decodeGetAllRequest, encodeGetAllResponse, options...),
		update:       grpc.NewServer(endpoints.UpdateUserEndpoint, decodeUpdateUserRequest, encodeResponse, options...),
		delete:       grpc.NewServer(endpoints.DeleteUserEndpoint, decodeUserIdRequest, encodeResponse, options...),
	}

	return server
}

func (srv *grpcUserServer) Authenticate(ctx context.Context, auth *pb.AuthRequest) (*pb.Response, error) {
	_, grpcResponse, err := srv.authenticate.ServeGRPC(ctx, auth)
	if err != nil {
		return &pb.Response{Code: 401}, err
	}
	return grpcResponse.(*pb.Response), err
}

func (srv *grpcUserServer) Create(ctx context.Context, user *pb.UserRequest) (*pb.Response, error) {
	_, grpcResponse, err := srv.createUser.ServeGRPC(ctx, user)
	if err != nil {
		return &pb.Response{Code: 401}, err
	}
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
	if err != nil {
		return &pb.Response{Code: 401}, err
	}
	return grpcResponse.(*pb.Response), err
}

func (srv *grpcUserServer) Delete(ctx context.Context, uid *pb.UserIDRequest) (*pb.Response, error) {
	_, grpcResponse, err := srv.delete.ServeGRPC(ctx, uid)
	if err != nil {
		return &pb.Response{Code: 401}, err
	}
	return grpcResponse.(*pb.Response), err
}

func (srv *grpcUserServer) GetAll(ctx context.Context, void *pb.Void) (*pb.UserColletionResponse, error) {
	_, grpcResponse, err := srv.getAll.ServeGRPC(ctx, void)

	return grpcResponse.(*pb.UserColletionResponse), err
}
