package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/caarlos0/env/v6"
	logkit "github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	grpcimp "github.com/neidersalgado/go-grpc-api/cmd/grpc-server/users"
	"github.com/neidersalgado/go-grpc-api/cmd/grpc-server/users/proto"
	pb "github.com/neidersalgado/go-grpc-api/cmd/grpc-server/users/proto"
	"github.com/neidersalgado/go-grpc-api/pkg/repository"
	"github.com/neidersalgado/go-grpc-api/pkg/users"
)

func main() {
	var logger logkit.Logger
	{
		logger = logkit.NewLogfmtLogger(os.Stderr)
		logger = logkit.With(logger, "ts", logkit.DefaultTimestampUTC)
		logger = logkit.With(logger, "caller", logkit.DefaultCaller)
	}

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	ls, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		panic(fmt.Sprintf("Could not create the listener %v", err))
	}

	repo, err := repository.NewMySQLUserRepository(logger)

	if err != nil {
		panic(fmt.Sprintf("mysql connection failed: %s", err.Error()))
	}

	userService := users.NewUserService(repo, logger)
	endpoints := grpcimp.NewGrpcUserServerEndpoints(userService)
	grpcUserServer := grpcimp.NewGrpcUserServer(*endpoints, logger)
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	proto.RegisterUsersServer(baseServer, grpcUserServer)

	go func() {
		log.Fatalln(baseServer.Serve(ls))
	}()

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = pb.RegisterUsersHandlerServer(context.Background(), gwmux, grpcUserServer)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://127.0.0.1:8090")
	log.Fatalln(gwServer.ListenAndServe())

}

type config struct {
	Port int    `env:"GRPCSERVICE_PORT" envDefault:"9000"`
	Host string `env:"GRPCSERVICE_HOST" envDefault:"127.0.0.1"`
}
