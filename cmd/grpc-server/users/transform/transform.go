package transform

import (
	pb "github.com/neidersalgado/go-grpc-api/cmd/grpc-server/users/proto"
	"github.com/neidersalgado/go-grpc-api/pkg/users"
)

func FromRequestToDomainS(userToMap pb.UserRequest) users.User {
	return users.User{
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
		PwdHash:               userToMap.PwdHash,
		Email:                 userToMap.Email,
		Name:                  userToMap.Name,
		Age:                   userToMap.Age,
		AdditionalInformation: userToMap.AdditionalInformation,
	}
}
