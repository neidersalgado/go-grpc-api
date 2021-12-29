package user

import (
	"github.com/neidersalgado/go-grpc-api/cmd/rest-service/user/pb"
	"github.com/neidersalgado/go-grpc-api/pkg/entities"
)

func transformUserEntityToRequest(userEntity entities.User) pb.UserRequest {
	return pb.UserRequest{
		UserId:                userEntity.UserID,
		Email:                 userEntity.Email,
		Name:                  userEntity.Name,
		Age:                   userEntity.Age,
		AdditionalInformation: userEntity.AdditionalInformation,
		Parents:               transformUserEntitiesToRequests(userEntity.Parent),
	}
}

func transformUserEntitiesToRequests(userEntities []entities.User) []*pb.UserRequest {
	if len(userEntities) == 0 {
		return []*pb.UserRequest{}
	}

	parentsRequest := make([]*pb.UserRequest, len(userEntities))

	for _, entity := range userEntities {
		parent := transformUserEntityToRequest(entity)
		parentsRequest = append(parentsRequest, &parent)
	}

	return parentsRequest
}

func transformUserIdToUserIdRequest(id string) *pb.UserIDRequest {
	return &pb.UserIDRequest{Email: id}
}

func transformUserRequestToEntity(userRequest pb.UserRequest) entities.User {
	return entities.User{
		UserID:                userRequest.UserId,
		Email:                 userRequest.Email,
		Name:                  userRequest.Name,
		Age:                   userRequest.Age,
		AdditionalInformation: userRequest.AdditionalInformation,
		Parent:                transformUserRequestsToEntities(userRequest.Parents),
	}
}

func transformUserResponseToEntity(userResponse pb.UserResponse) entities.User {
	return entities.User{
		UserID:                userResponse.UserId,
		Email:                 userResponse.Email,
		Name:                  userResponse.Name,
		Age:                   userResponse.Age,
		AdditionalInformation: userResponse.AdditionalInformation,
		Parent:                transformUserResponsesToEntities(userResponse.Parents),
	}
}

func transformUserRequestsToEntities(userRequest []*pb.UserRequest) []entities.User {
	if len(userRequest) == 0 {
		return []entities.User{}
	}

	parentsEntities := make([]entities.User, len(userRequest))

	for _, entity := range userRequest {
		parent := transformUserRequestToEntity(*entity)
		parentsEntities = append(parentsEntities, parent)
	}

	return parentsEntities
}

func transformUserResponsesToEntities(userResponses []*pb.UserResponse) []entities.User {
	if len(userResponses) == 0 {
		return []entities.User{}
	}

	entities := make([]entities.User, len(userResponses))

	for index, entity := range userResponses {
		parent := transformUserResponseToEntity(*entity)
		entities[index] = parent
	}

	return entities
}
