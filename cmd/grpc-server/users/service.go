package grpc

import (
	"context"

	"github.com/neidersalgado/go-camp-grpc/cmd/grpc-server/pb"
	"github.com/neidersalgado/go-camp-grpc/pkg/users"
)

type GrpcService struct {
	userService users.Service
}

func NewGrpcUserService(us users.Service) *GrpcService {
	return &GrpcService{
		userService: us,
	}
}

func (s *GrpcService) Authenticate(ctx context.Context, auth *pb.AuthRequest) (*pb.Response, error) {
	return &pb.Response{}, nil
}

func (s *GrpcService) Create(ctx context.Context, usr *pb.UserRequest) (*pb.Response, error) {
	userToCreate := users.User{
		UserId:                usr.UserId,
		Email:                 usr.Email,
		PwdHash:               usr.PwdHash,
		Name:                  usr.Name,
		Age:                   usr.Age,
		AdditionalInformation: usr.AdditionalInformation,
	}

	err := s.userService.Create(ctx, userToCreate)

	if err != nil {
		return &pb.Response{Code: 400}, err
	}

	return &pb.Response{Code: 200}, nil
}

func (s *GrpcService) Get(ctx context.Context, usrID *pb.UserIDRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{}, nil
}

func (s *GrpcService) Update(ctx context.Context, usr *pb.UserRequest) (*pb.Response, error) {
	return &pb.Response{}, nil
}

func (s *GrpcService) Delete(ctx context.Context, usrID *pb.UserIDRequest) (*pb.Response, error) {
	return &pb.Response{}, nil
}

func (s *GrpcService) GetAll(ctx context.Context, void *pb.Void) (*pb.UserColletionResponse, error) {
	return &pb.UserColletionResponse{}, nil
}
