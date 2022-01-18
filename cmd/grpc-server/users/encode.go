package grpc

import (
	"context"
	"errors"

	pb "github.com/neidersalgado/go-grpc-api/cmd/grpc-server/users/proto"
)

const (
	INVALIDENCODEDATA = "invalid input data encode"
	ALREADYEXIST      = "user already exists"
)

func encodeCreateUserResponse(ctx context.Context, resp interface{}) (interface{}, error) {

	respData, validCast := resp.(createUserResponse)
	if !validCast {
		return nil, errors.New(INVALIDENCODEDATA)
	}

	if respData.Error != nil {
		if respData.Error.Error() == ALREADYEXIST {
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
		return nil, errors.New(INVALIDENCODEDATA)
	}

	return &pb.UserResponse{
		PwdHash:               usrResponse.PwdHash,
		Email:                 usrResponse.Email,
		Name:                  usrResponse.Name,
		Age:                   usrResponse.Age,
		AdditionalInformation: usrResponse.AdditionalInformation,
	}, nil
}

func encodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	response, validCast := resp.(Response)
	if !validCast {
		return &pb.Response{Code: 400}, errors.New(INVALIDENCODEDATA)
	}
	return &pb.Response{Code: pb.Response_CodeResult(response.Code)}, nil
}

func encodeGetAllResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	usrs, validCast := resp.(getAllUsersResponse)
	if !validCast {
		return &pb.Response{Code: 400}, errors.New(INVALIDENCODEDATA)
	}

	usersResponse := make([]*pb.UserResponse, len(usrs.Users))

	for index, usr := range usrs.Users {
		userPb := pb.UserResponse{
			PwdHash:               usr.PwdHash,
			Email:                 usr.Email,
			Name:                  usr.Name,
			Age:                   usr.Age,
			AdditionalInformation: usr.AdditionalInformation,
		}
		usersResponse[index] = &userPb
	}
	usersCollection := &pb.UserColletionResponse{
		Users: usersResponse,
	}
	return usersCollection, nil
}
