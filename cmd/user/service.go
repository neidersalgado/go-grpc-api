package user

import (
	"context"
	"fmt"

	"github.com/neidersalgado/go-camp-grpc/cmd/user/pb"
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
	return &pb.Response{Code: 501}, nil
}

func (us *UsersService) Create(ctx context.Context, userReq *pb.UserRequest) (*pb.Response, error) {
	if err := us.repository.Create(*userReq); err != nil {
		fmt.Printf(
			fmt.Sprintf("\n couldn't create user With ID: %d, \n Error: %v \n", userReq.UserId, err.Error()),
		)
		return &pb.Response{Code: 400}, fmt.Errorf("couldn't create user With ID: %d, \n Error: %v", userReq.UserId, err.Error())
	}
	return &pb.Response{Code: 200}, nil
}

func (us *UsersService) Get(ctx context.Context, ID *pb.UserID) (*pb.UserResponse, error) {
	user, err := us.repository.Get(ID.ID)

	if err != nil {
		return &pb.UserResponse{}, fmt.Errorf("couldn't get user With id: %s, \n Error: %v", ID.GetID(), err.Error())
	}

	return &user, nil
}

func (us *UsersService) Update(ctx context.Context, user *pb.UserRequest) (*pb.Response, error) {
	if err := us.repository.Update(*user); err != nil {
		return &pb.Response{Code: 400}, fmt.Errorf("couldn't update user with id: %v, \n error: %v", user.UserId, err.Error())
	}
	return &pb.Response{Code: 200}, nil
}

func (us *UsersService) Delete(ctx context.Context, ID *pb.UserID) (*pb.Response, error) {
	if err := us.repository.Delete(ID.ID); err != nil {
		return &pb.Response{Code: 400}, fmt.Errorf("couldn't delete user With ID: %s, \n Error: %v", ID.ID, err.Error())
	}

	return &pb.Response{Code: 200}, nil
}

func (us *UsersService) GetAll(context.Context, *pb.Void) (*pb.UserColletionResponse, error) {
	users, err := us.repository.GetAll()
	list := pb.UserColletionResponse{Users: users.Users}
	if err != nil {
		return &pb.UserColletionResponse{}, fmt.Errorf("couldn't get all users \n Error: %v", err.Error())
	}
	return &list, nil
}
