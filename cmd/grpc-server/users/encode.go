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

func encodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	response, validCast := resp.(Response)
	if !validCast {
		return &pb.Response{Code: 400}, errors.New("invalid input data to encode")
	}
	return &pb.Response{Code: pb.Response_CodeResult(response.Code)}, nil
}

func encodeGetAllResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	usrs, validCast := resp.(getAllUsersResponse)
	if !validCast {
		return &pb.Response{Code: 400}, errors.New("invalid input data to encode")
	}

	usersResponse := make([]*pb.UserResponse, len(usrs.Users))

	for _, usr := range usrs.Users {
		userPb := pb.UserResponse{
			UserId:                usr.UserId,
			PwdHash:               usr.PwdHash,
			Email:                 usr.Email,
			Name:                  usr.Name,
			Age:                   usr.Age,
			AdditionalInformation: usr.AdditionalInformation,
		}
		usersResponse = append(usersResponse, &userPb)
	}
	usersCollection := &pb.UserColletionResponse{
		Users: usersResponse,
	}
	return usersCollection, nil
}
