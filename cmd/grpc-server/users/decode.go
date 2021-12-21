package grpc

import (
	"context"
	"errors"

	"github.com/neidersalgado/go-camp-grpc/cmd/grpc-server/pb"
	domain "github.com/neidersalgado/go-camp-grpc/pkg/users"
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

func decodeUserIdRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*pb.UserIDRequest)
	if !validCast {
		return nil, errors.New("invalid input data decode")
	}

	return userIdRequest{Email: reqData.Email}, nil
}

func decodeAuthUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*pb.AuthRequest)

	if !validCast {
		return domain.Auth{}, errors.New("invalid input data decode")
	}
	return domain.Auth{Mail: reqData.Name, Hash: reqData.Hash}, nil
}

func decodeGetAllRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	_, validCast := grpcReq.(*pb.Void)

	if !validCast {
		return void{}, errors.New("invalid input data decode")
	}

	return void{}, nil
}

func decodeUpdateUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	updateData, validCast := grpcReq.(*pb.UserRequest)

	if !validCast {
		return nil, errors.New("invalid input data decode")
	}
	usr := UserRequest{
		UserId:                updateData.UserId,
		PwdHash:               updateData.PwdHash,
		Email:                 updateData.Email,
		Name:                  updateData.Name,
		Age:                   updateData.Age,
		AdditionalInformation: updateData.AdditionalInformation,
	}
	return updateUserRequest{UserRequest: usr}, nil
}
