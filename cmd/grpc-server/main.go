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

	ls, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GrpcHost, cfg.GrpcPort))
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
		Addr:    fmt.Sprintf("%s:%d", cfg.HttpHost, cfg.HttpPort),
		Handler: gwmux,
	}

	log.Println(fmt.Sprintf("Serving gRPC-Gateway on http://%s:%d", cfg.HttpHost, cfg.HttpPort))
	log.Fatalln(gwServer.ListenAndServe())

}

type config struct {
	GrpcPort int    `env:"GRPC_SERVICE_PORT" envDefault:"9000"`
	GrpcHost string `env:"GRPC_SERVICE_HOST" envDefault:"127.0.0.1"`
	HttpPort int    `env:"HTTP_SERVICE_PORT" envDefault:"8090"`
	HttpHost string `env:"HTTP_SERVICE_HOST" envDefault:"127.0.0.1"`
}
