package transform

import (
	"github.com/neidersalgado/go-grpc-api/cmd/grpc-server/users/pb"
	"github.com/neidersalgado/go-grpc-api/pkg/users"
)

func FromRequestToDomainS(userToMap pb.UserRequest) users.User {
	return users.User{
		UserId:                userToMap.UserId,
		PwdHash:               userToMap.PwdHash,
		Email:                 userToMap.Email,
		Name:                  userToMap.Name,
		Age:                   userToMap.Age,
		AdditionalInformation: userToMap.AdditionalInformation,
	}
}

//ToGrpcUser maps a domain user to a grpc user
func FromDomainToResponseS(userToMap users.User) pb.UserResponse {
	return pb.UserResponse{
		UserId:                userToMap.UserId,
		PwdHash:               userToMap.PwdHash,
		Email:                 userToMap.Email,
		Name:                  userToMap.Name,
		Age:                   userToMap.Age,
		AdditionalInformation: userToMap.AdditionalInformation,
	}
}
