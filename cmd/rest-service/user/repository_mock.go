package user

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/neidersalgado/go-grpc-api/pkg/entities"
)

type MockRepository struct {
	mock.Mock
}

func (rm *MockRepository) Authenticate(email string, hash string) (bool, error) {
	args := rm.Called(email, hash)
	return args.Bool(0), args.Error(1)
}

func (rm *MockRepository) Create(user entities.User) error {
	args := rm.Called(user)
	return args.Error(0)
}

func (rm *MockRepository) Update(ctx context.Context, user entities.User) error {
	args := rm.Called(ctx, user)
	return args.Error(0)
}

func (rm *MockRepository) Get(ctx context.Context, userID string) (entities.User, error) {
	args := rm.Called(ctx, userID)
	return args.Get(0).(entities.User), args.Error(0)
}

func (rm *MockRepository) List(ctx context.Context) ([]entities.User, error) {
	args := rm.Called(ctx)
	return args.Get(0).([]entities.User), args.Error(0)
}

func (rm *MockRepository) Delete(ctx context.Context, userID string) error {
	args := rm.Called(ctx, userID)
	return args.Error(0)
}
