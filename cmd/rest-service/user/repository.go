package user

import (
	"context"
	"errors"
	"fmt"
	logSys "log"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"

	"github.com/neidersalgado/go-grpc-api/cmd/rest-service/user/pb"
	"github.com/neidersalgado/go-grpc-api/pkg/entities"
)

type config struct {
	Port int    `env:"GRPCSERVICE_PORT" envDefault:"9000"`
	//Host string `env:"GRPCSERVICE_HOST" envDefault:"127.0.0.1"`
	Host string `env:"GRPCSERVICE_HOST" envDefault:"172.19.0.3"`
}

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
	serverCon, err := OpenServerConnection(context.Background(), up.logger)

	if err != nil {
		up.logger.Log("proxy", fmt.Sprintf("did not connect to server: %s", err))
	}
	up.logger.Log("Proxy", "create", "creation", "Connection Ok")
	defer serverCon.dispose()
	c := serverCon.client
	externalUser := transformUserEntityToRequest(user)
	up.logger.Log("Proxy", "create", "transform", "completed")
	result, errorFromCall := c.Create(serverCon.context, &externalUser)
	if errorFromCall != nil {
		up.logger.Log("client", "create", "create", fmt.Sprintf("Error: %+v", errorFromCall.Error()))
		return errors.New(fmt.Sprintf("Error Creating User  error: %v", errorFromCall.Error()))
	}

	if result.GetCode() != http.StatusOK {
		up.logger.Log("client", "create", "code", fmt.Sprintf("status: %+v", result.Code))
		return errors.New(fmt.Sprintf("Error Creating User  error: %v, response code: %v", errorFromCall.Error(), result.Code))
	}

	return nil
}

func (up ProxyRepository) Update(ctx context.Context, user entities.User) error {
	up.logger.Log("proxy", "update", "Updating", fmt.Sprintf("user: %+v", user))
	serverCon, err := OpenServerConnection(ctx, up.logger)

	if err != nil {
		up.logger.Log("proxy", fmt.Sprintf("did not connect to server: %s", err))
		return err
	}

	defer serverCon.dispose()
	c := serverCon.client
	updateRequest := transformUserEntityToRequest(user)
	response, errorFromCall := c.Update(ctx, &updateRequest)
	up.logger.Log("client", "update", "Update", fmt.Sprintf("Request: %+v", updateRequest))

	if errorFromCall != nil {
		up.logger.Log("client", "update", "create", fmt.Sprintf("Error: %+v", errorFromCall.Error()))
		return errors.New(fmt.Sprintf("Error Updating User  error: %v, Response:%v\n", errorFromCall.Error(), response))
	}
	up.logger.Log("client", "update", "Update", fmt.Sprintf("respose: %+v", response))
	return nil
}

func (up ProxyRepository) Get(ctx context.Context, userID string) (entities.User, error) {
	fmt.Printf("repository.Get with id: %v.\n", userID)
	serverCon, err := OpenServerConnection(ctx, up.logger)

	if err != nil {
		up.logger.Log("proxy", fmt.Sprintf("did not connect to server: %s", err.Error()))
	}

	defer serverCon.dispose()
	c := serverCon.client
	userIDpb := transformUserIdToUserIdRequest(userID)
	userResponse, errorFromCall := c.Get(ctx, userIDpb)
	if errorFromCall != nil {
		up.logger.Log("error", fmt.Sprintf("Error Getting User: %v", errorFromCall.Error()))
		fmt.Printf("service. Error:%v \n", errorFromCall.Error())
		return entities.User{}, errors.New(fmt.Sprintf("Error Getting User  error: %v", errorFromCall.Error()))
	}
	userEntity := transformUserResponseToEntity(*userResponse)

	return userEntity, nil
}

func (up ProxyRepository) List(ctx context.Context) ([]entities.User, error) {
	fmt.Printf("repository.List \n")
	serverCon, err := OpenServerConnection(ctx, up.logger)

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
	serverCon, err := OpenServerConnection(ctx, up.logger)

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

func OpenServerConnection(ctx context.Context, logger log.Logger) (*ServerConnection, error) {

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	logger.Log("********** Dial", fmt.Sprintf("%s:%s", cfg.Host, strconv.Itoa(cfg.Port)))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Host, strconv.Itoa(cfg.Port)), grpc.WithInsecure())

	if err != nil {
		logSys.Fatalf("did not connect to server: %s", err)
		return nil, err
	}

	ctxTO, cancel := context.WithTimeout(ctx, 10*time.Second)

	c := pb.NewUsersClient(conn)

	return &ServerConnection{
		client:  c,
		context: ctxTO,
		dispose: func() {
			cancel()
			conn.Close()
		},
	}, nil

}

type ServerConnection struct {
	client  pb.UsersClient
	context context.Context
	dispose func()
}
