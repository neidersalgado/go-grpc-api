package grpc

import (
	"context"
	"errors"

	"github.com/neidersalgado/go-grpc-api/cmd/grpc-server/users/proto"
	domain "github.com/neidersalgado/go-grpc-api/pkg/users"
)

const (
	INVALIDADECODEDATA = "invalid input data decode"
)

func decodeCreateUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*proto.UserRequest)
	if !validCast {
		return nil, errors.New(INVALIDADECODEDATA)
	}
	usr := UserRequest{
		PwdHash:               reqData.PwdHash,
		Email:                 reqData.Email,
		Name:                  reqData.Name,
		Age:                   reqData.Age,
		AdditionalInformation: reqData.AdditionalInformation,
	}

	return createUserRequest{UserRequest: usr}, nil
}

func decodeUserIdRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*proto.UserIDRequest)
	if !validCast {
		return nil, errors.New(INVALIDADECODEDATA)
	}

	return userIdRequest{Email: reqData.Email}, nil
}

func decodeAuthUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	reqData, validCast := grpcReq.(*proto.AuthRequest)

	if !validCast {
		return domain.Auth{}, errors.New(INVALIDADECODEDATA)
	}
	return domain.Auth{Mail: reqData.Name, Hash: reqData.Hash}, nil
}

func decodeGetAllRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	_, validCast := grpcReq.(*proto.Void)

	if !validCast {
		return void{}, errors.New(INVALIDADECODEDATA)
	}

	return void{}, nil
}

func decodeUpdateUserRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	updateData, validCast := grpcReq.(*proto.UserRequest)

	if !validCast {
		return nil, errors.New(INVALIDADECODEDATA)
	}
	usr := UserRequest{
		PwdHash:               updateData.PwdHash,
		Email:                 updateData.Email,
		Name:                  updateData.Name,
		Age:                   updateData.Age,
		AdditionalInformation: updateData.AdditionalInformation,
	}
	return updateUserRequest{UserRequest: usr}, nil
}
