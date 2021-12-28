package users

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (mock *repositoryMock) Create(ctx context.Context, user User) error {
	args := mock.Called(ctx, user)
	return args.Error(0)
}

func (mock *repositoryMock) GetByEmail(ctx context.Context, email string) (User, error) {
	args := mock.Called(ctx, email)
	return args.Get(0).(User), args.Error(1)
}

func (mock *repositoryMock) Update(ctx context.Context, user User) error {
	args := mock.Called(ctx, user)
	return args.Error(0)
}

func (mock *repositoryMock) Delete(ctx context.Context, email string) error {
	args := mock.Called(ctx, email)
	return args.Error(0)
}

func (mock *repositoryMock) GetAll(ctx context.Context) ([]User, error) {
	args := mock.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}


