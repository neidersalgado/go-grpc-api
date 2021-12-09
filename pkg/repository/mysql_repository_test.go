package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/neidersalgado/go-camp-grpc/pkg/users"
)

var repository *MySQLUserRepository

func TestMain(m *testing.M) {
	var err error
	repository, err = NewMySQLUserRepository()

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func Test_Create_ValidData_ReturnNil(t *testing.T) {
	var err error
	userToCreate := users.User{Email: "mail@mail.com", Name: "name", Age: 25, AdditionalInformation: "none"}

	err = repository.Create(context.Background(), userToCreate)

	assert.Nil(t, err)
}
