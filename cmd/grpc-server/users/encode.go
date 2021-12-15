package grpc

import (
	"context"
	"errors"

	"github.com/neidersalgado/go-camp-grpc/cmd/grpc-server/pb"
)

func encodeCreateUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {

	respData, validCast := resp.(createUserResponse)
	if !validCast {
		return nil, errors.New("invalid input data encode")
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

func encodeGetUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {

	usrResponse, validCast := resp.(getUserResponse)

	if !validCast {
		return nil, errors.New("invalid input data to encode")
	}

	return &pb.UserResponse{
		UserId:                usrResponse.UserId,
		PwdHash:               usrResponse.PwdHash,
		Email:                 usrResponse.Email,
		Name:                  usrResponse.Name,
		Age:                   usrResponse.Age,
		AdditionalInformation: usrResponse.AdditionalInformation,
	}, nil
}
