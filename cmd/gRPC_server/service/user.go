package userservice

import (
	"context"
	"fmt"

	"github.com/neidersalgado/go-camp-grpc/cmd/gRPC_server/pb"
)

type Repository interface {
	Create(user pb.UserRequest) error
	Get(userID string) (pb.UserResponse, error)
	Update(user pb.UserRequest) error
	Delete(userID string) error
	GetAll() (pb.UserColletionResponse, error)
}

type UsersService struct {
	pb.UnimplementedUsersServer
	repository Repository
}

func NewUserService(repo Repository) *UsersService {
	return &UsersService{
		repository: repo,
	}
}

func (us *UsersService) Authenticate(context.Context, *pb.AuthRequest) (*pb.Response, error) {
	return &pb.Response{Code: 501}, fmt.Errorf("method Authenticate not implemented")
}

func (us *UsersService) Create(ctx context.Context, userReq *pb.UserRequest) (*pb.Response, error) {
	if err := us.repository.Create(*userReq); err != nil {
		return &pb.Response{Code: 500}, fmt.Errorf("couldn't create user With ID: %s, \n Error: %v", userReq.Id, err)
	}

	return &pb.Response{Code: 200}, nil
}

func (us *UsersService) Get(ctx context.Context, ID *pb.UserID) (*pb.UserResponse, error) {
	user, err := us.repository.Get(ID.ID)

	if err != nil {
		return &pb.UserResponse{}, fmt.Errorf("couldn't get user With ID: %s, \n Error: %v", ID.ID, err)
	}

	return &user, nil
}

func (us *UsersService) Update(context.Context, *pb.UserRequest) (*pb.Response, error) {
	return &pb.Response{}, fmt.Errorf("method update not implemented")
}

func (us *UsersService) Delete(ctx context.Context, ID *pb.UserID) (*pb.Response, error) {
	if err := us.repository.Delete(ID.ID); err != nil {
		return &pb.Response{Code: 500}, fmt.Errorf("couldn't delete user With ID: %s, \n Error: %v", ID.ID, err)
	}

	return &pb.Response{Code: 200}, nil
}

func (us *UsersService) GetAll(context.Context, *pb.Void) (*pb.UserColletionResponse, error) {
	return nil, fmt.Errorf("method getAll not implemented")
}
