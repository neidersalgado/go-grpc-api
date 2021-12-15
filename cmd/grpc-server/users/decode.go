package grpc

import (
	"context"
	"errors"

	"github.com/neidersalgado/go-camp-grpc/cmd/grpc-server/pb"
)

func decodeCreateUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*pb.UserRequest)
	if !validCast {
		return nil, errors.New("invalid input data decode")
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

func decodeGetUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*pb.UserIDRequest)
	if !validCast {
		return nil, errors.New("invalid input data decode")
	}

	return getUserRequest{Email: reqData.Email}, nil
}
