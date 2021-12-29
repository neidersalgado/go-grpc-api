package transform

import "github.com/neidersalgado/go-grpc-api/cmd/grpc-server/users/pb"

func FromRequestToDomainS(userToMap pb.UserRequest) domain.User {
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
func FromDomainToResponseS(userToMap domain.User) pb.UserResponse {
	return pb.UserResponse{
		UserId:                userToMap.UserId,
		PwdHash:               userToMap.PwdHash,
		Email:                 userToMap.Email,
		Name:                  userToMap.Name,
		Age:                   userToMap.Age,
		AdditionalInformation: userToMap.AdditionalInformation,
	}
}
