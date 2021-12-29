package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"

	"github.com/neidersalgado/go-grpc-api/cmd/rest-service/user/pb"
	"github.com/neidersalgado/go-grpc-api/pkg/entities"
)

type ProxyRepository struct {
	logger log.Logger
}

func NewProxyRepository(log log.Logger) *ProxyRepository {
	return &ProxyRepository{
		logger: log,
	}
}

func (up ProxyRepository) Authenticate(email string, hash string) (bool, error) {
	fmt.Println("User well authenticated")
	return true, nil
}

func (up ProxyRepository) Create(user entities.User) error {
	serverCon, err := OpenServerConnection()

	if err != nil {
		up.logger.Log("proxy", fmt.Sprintf("did not connect to server: %s", err))
	}

	defer serverCon.dispose()
	c := serverCon.client
	externalUser := transformUserEntityToRequest(user)
	result, errorFromCall := c.Create(serverCon.context, &externalUser)
	if errorFromCall != nil {
		return errors.New(fmt.Sprintf("Error Creating User  error: %v", errorFromCall.Error()))
	}

	if result.GetCode() != http.StatusOK {
		return errors.New(fmt.Sprintf("Error Creating User  error: %v, response code: %v", errorFromCall, result.Code))
	}

	return nil
}

func (up ProxyRepository) Update(ctx context.Context, user entities.User) error {
	fmt.Printf("repository.Update user user id: %v \n", user.UserID)
	serverCon, err := OpenServerConnection()

	if err != nil {
		up.logger.Log("proxy", fmt.Sprintf("did not connect to server: %s", err))
	}

	defer serverCon.dispose()
	c := serverCon.client
	updateRequest := transformUserEntityToRequest(user)
	response, errorFromCall := c.Update(ctx, &updateRequest)
	fmt.Printf("service.repo.Update")

	if errorFromCall != nil {
		return errors.New(fmt.Sprintf("Error Updating User  error: %v, Response:%v\n", errorFromCall.Error(), response))
	}
	return nil
}

func (up ProxyRepository) Get(ctx context.Context, userID string) (entities.User, error) {
	fmt.Printf("repository.Get with id: %v.\n", userID)
	serverCon, err := OpenServerConnection()

	if err != nil {
		up.logger.Log("proxy", fmt.Sprintf("did not connect to server: %s", err))
	}

	defer serverCon.dispose()
	c := serverCon.client
	userIDpb := transformUserIdToUserIdRequest(userID)
	fmt.Printf("repository.Getcasting user to  call grpc service \n")
	userResponse, errorFromCall := c.Get(ctx, userIDpb)
	fmt.Printf("service. GEt user sesponse :%v \n", userResponse)
	if errorFromCall != nil {
		return entities.User{}, errors.New(fmt.Sprintf("Error Creating User  error: %v", errorFromCall.Error()))
	}

	userEntity := transformUserResponseToEntity(*userResponse)

	return userEntity, nil
}

func (up ProxyRepository) List(ctx context.Context) ([]entities.User, error) {
	fmt.Printf("repository.List \n")
	serverCon, err := OpenServerConnection()

	if err != nil {
		up.logger.Log("proxy", fmt.Sprintf("did not connect to server: %s", err))
	}
	void := pb.Void{}
	defer serverCon.dispose()
	c := serverCon.client
	usersResponse, errorFromCall := c.GetAll(ctx, &void)
	fmt.Printf("service.repo.list")
	if errorFromCall != nil {
		return []entities.User{}, errors.New(fmt.Sprintf("Error list Users error: %v", errorFromCall.Error()))
	}
	userEntities := transformUserResponsesToEntities(usersResponse.Users)
	return userEntities, nil
}

func (up ProxyRepository) Delete(ctx context.Context, userID string) error {
	fmt.Printf("repository.Delete user user id: %v \n", userID)
	serverCon, err := OpenServerConnection()

	if err != nil {
		up.logger.Log("proxy", fmt.Sprintf("did not connect to server: %s", err))
	}

	defer serverCon.dispose()
	c := serverCon.client
	deleteRequest := transformUserIdToUserIdRequest(userID)
	response, errorFromCall := c.Delete(ctx, deleteRequest)
	fmt.Printf("service.repo.Delete")

	if errorFromCall != nil {
		return errors.New(fmt.Sprintf("Error deleting User  error: %v, Response:%v\n", errorFromCall.Error(), response))
	}
	return nil
}

func OpenServerConnection() (*ServerConnection, error) {

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		return nil, err //unreached?
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	c := pb.NewUsersClient(conn)

	return &ServerConnection{client: c, context: ctx, dispose: func() {
		cancel()
		conn.Close()

	}}, nil

}

type ServerConnection struct {
	client  pb.UsersClient
	context context.Context
	dispose func()
}
