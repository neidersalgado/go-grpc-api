package users

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Create_ValidData_OkResult(t *testing.T) {
	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToAdd := User{Email: "test@gmail.com", Name: "John", Age: 24, AdditionalInformation: "none"}
	repository.On("Create", context.Background(), userToAdd).Return(nil)
	repository.On("GetByEmail", context.Background(), userToAdd.Email).Return(User{}, nil)

	err := service.Create(context.Background(), userToAdd)

	assert.Nil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "Create", 1)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}

func Test_Create_DuplicateData_ReturnErrorGet(t *testing.T) {

	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToAdd := User{Email: "test@gmail.com", Name: "John", Age: 24, AdditionalInformation: "none"}
	repository.On("GetByEmail", context.Background(), userToAdd.Email).Return(User{}, errors.New(""))

	err := service.Create(context.Background(), userToAdd)

	assert.NotNil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}

func Test_Create_RepositoryGetError_ReturnError(t *testing.T) {

	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToAdd := User{Email: "test@gmail.com", Name: "John", Age: 24, AdditionalInformation: "none"}
	repository.On("GetByEmail", context.Background(), userToAdd.Email).Return(User{}, errors.New(""))

	err := service.Create(context.Background(), userToAdd)

	assert.NotNil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}

func Test_Create_DuplicateData_ReturnError(t *testing.T) {

	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToAdd := User{UserId: 20, Email: "test@gmail.com", Name: "John", Age: 24, AdditionalInformation: "none"}
	repository.On("GetByEmail", context.Background(), userToAdd.Email).Return(userToAdd, nil)

	err := service.Create(context.Background(), userToAdd)

	assert.NotNil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}

func Test_Create_CreateFails_ReturnErrorCreate(t *testing.T) {

	repository := repositoryMock{}
	service := NewUserService(&repository)
	userToAdd := User{Email: "test@gmail.com", Name: "John", Age: 24, AdditionalInformation: "none"}
	repository.On("Create", context.Background(), userToAdd).Return(errors.New(""))
	repository.On("GetByEmail", context.Background(), userToAdd.Email).Return(User{}, nil)

	err := service.Create(context.Background(), userToAdd)

	assert.NotNil(t, err)
	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "Create", 1)
	repository.AssertNumberOfCalls(t, "GetByEmail", 1)
}
