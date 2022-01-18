package repository

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/neidersalgado/go-grpc-api/pkg/users"
)

var repository *MySQLUserRepository

func TestMain(m *testing.M) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	var err error
	repository, err = NewMySQLUserRepository(logger)

	if err != nil {
		panic(err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}
/*
func Test_Create_ValidData_ReturnNil(t *testing.T) {
	var err error
	userToCreate := users.User{Email: "mail@mail.com", Name: "name", Age: 25, AdditionalInformation: "none"}

	err = repository.Create(context.Background(), userToCreate)

	assert.Nil(t, err)
}
//*/s