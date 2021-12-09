package main

import (
	"fmt"
	"net"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"github.com/neidersalgado/go-camp-grpc/cmd/grpc-server/pb"
	grpcImp "github.com/neidersalgado/go-camp-grpc/cmd/grpc-server/users"
	"github.com/neidersalgado/go-camp-grpc/pkg/repository"
	"github.com/neidersalgado/go-camp-grpc/pkg/users"
)

func main() {

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))

	if err != nil {
		panic(fmt.Sprintf("Could not create the listener %v", err))
	}

	repository, err := repository.NewMySQLUserRepository()

	if err != nil {
		panic(fmt.Sprintf("mysql connection failed: %s", err))
	}

	userService := users.NewUserService(repository)
	endpoints := grpcImp.NewGrpcUserServerEndpoints(*userService)
	grpcUserServer := grpcImp.NewGrpcUserServer(*endpoints, logger)
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	pb.RegisterUsersServer(baseServer, grpcUserServer)

	if err := baseServer.Serve(ls); err != nil {
		panic(fmt.Sprintf("failed to serve: %s", err))
	}
}

type config struct {
	Port int `env:"GRPCSERVICE_PORT" envDefault:"9000"`
}
