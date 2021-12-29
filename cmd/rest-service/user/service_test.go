package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/neidersalgado/go-grpc-api/pkg/entities"
)

func TestCasesCreate(t *testing.T) {
	ctx := context.Background()
	errRepo := errors.New("Not Implemented")

	var createUserCases []struct {
		testCaseName string
		data         entities.User
		err          error
		expectResult error
	} = []struct {
		testCaseName string
		data         entities.User
		err          error
		expectResult error
	}{
		{"user_well_created", getUser(), nil, nil},
		{"user_error_creating", getUser(), errRepo, errRepo},
	}
	for _, createCase := range createUserCases {
		repoMock := new(MockRepository)
		service := NewDefaultUserService(repoMock)
		repoMock.On("Create", createCase.data).Return(createCase.err)
		t.Run(createCase.testCaseName, func(t *testing.T) {
			result := service.CreateUser(ctx, getUser())
			if result != nil {
				assert.EqualError(t, createCase.expectResult, result.Error())
			}
		})
	}
}

func TestCasesUpdate(t *testing.T) {
	ctx := context.Background()
	errRepo := errors.New("Not Implemented")

	var updateUserCases []struct {
		testCaseName string
		data         entities.User
		err          error
		expectResult error
	} = []struct {
		testCaseName string
		data         entities.User
		err          error
		expectResult error
	}{
		{"user_well_created", getUser(), nil, nil},
		{"user_error_creating", getUser(), errRepo, errRepo},
	}
	for _, createCase := range updateUserCases {
		repoMock := new(MockRepository)
		service := NewDefaultUserService(repoMock)
		repoMock.On("Update", ctx, createCase.data).Return(createCase.err)
		t.Run(createCase.testCaseName, func(t *testing.T) {
			result := service.UpdateUser(ctx, getUser())
			if result != nil {
				assert.EqualError(t, createCase.expectResult, result.Error())
			}
		})
	}
}

func getUser() entities.User {
	return entities.User{
		Email:                 "mail@mail.com",
		Name:                  "fakeName",
		Age:                   23,
		AdditionalInformation: "none",
	}
}
