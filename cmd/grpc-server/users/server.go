package grpc

import (
	"context"
	"errors"
	"fmt"

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
		createUser:   grpc.NewServer(endpoints.CreateUserEndpoint, decodeCreateUserRequest, encodeCreateUserResponse, options...),
		getAll:       grpc.NewServer(endpoints.getAllEndpoint, nil, nil, options...),
		update:       grpc.NewServer(endpoints.UpdateUserEndpoint, nil, nil, options...),
		delete:       grpc.NewServer(endpoints.DeleteUserEndpoint, nil, nil, options...),
	}

	return server
}

func decodeCreateUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*pb.UserRequest)
	if !validCast {
		return nil, errors.New("invalid input data")
	}
	usr := UserRequest{
		UserId:                reqData.UserId,
		PwdHash:               reqData.PwdHash,
		Email:                 reqData.Email,
		Name:                  reqData.Name,
		Age:                   reqData.Age,
		AdditionalInformation: reqData.AdditionalInformation,
	}

	return createUserRequest{UserRequest: usr}, nil
}

func encodeCreateUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {

	fmt.Println("encode before cast", resp)
	respData, validCast := resp.(createUserResponse)
	fmt.Println("Cast", respData, validCast)
	if !validCast {
		return nil, errors.New("invalid input data")
	}

	if respData.Error != nil {
		if respData.Error.Error() == "user already exists" {
			return &pb.Response{Code: pb.Response_FAILED}, nil
		} else {
			return &pb.Response{Code: pb.Response_INVALIDINPUT}, nil
		}
	}

	return &pb.Response{Code: pb.Response_OK}, nil
}
