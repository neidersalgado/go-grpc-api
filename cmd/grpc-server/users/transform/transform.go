package transform

import (
	"github.com/neidersalgado/go-camp-grpc/cmd/grpc-server/pb"
	domain "github.com/neidersalgado/go-camp-grpc/pkg/users"
)

func FromRequestToDomain(userToMap pb.UserRequest) domain.User {
	return domain.User{
		UserId:                userToMap.UserId,
		PwdHash:               userToMap.PwdHash,
		Email:                 userToMap.Email,
		Name:                  userToMap.Name,
		Age:                   userToMap.Age,
		AdditionalInformation: userToMap.AdditionalInformation,
	}
}

//ToGrpcUser maps a domain user to a grpc user
func FromDomainToResponse(userToMap domain.User) pb.UserResponse {
	return pb.UserResponse{
		UserId:                userToMap.UserId,
		PwdHash:               userToMap.PwdHash,
		Email:                 userToMap.Email,
		Name:                  userToMap.Name,
		Age:                   userToMap.Age,
		AdditionalInformation: userToMap.AdditionalInformation,
	}
}
