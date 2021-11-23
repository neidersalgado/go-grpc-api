package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/caarlos0/env/v6"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"

	"github.com/neidersalgado/go-camp-grpc/cmd/user"
	"github.com/neidersalgado/go-camp-grpc/cmd/user/pb"
)

func main() {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	fmt.Println("Listen TCP")
	if err != nil {
		log.Fatalf("Could not create the listener %v", err)
	}

	db := dbConn()
	fmt.Println("DB Connect")
	defer db.Close()
	userRepo := user.NewMySQLUserRepository(db)
	server := grpc.NewServer()
	fmt.Println("Config Server")
	pb.RegisterUsersServer(server, user.NewUserService(userRepo))
	fmt.Println("Serving Service")
	if err := server.Serve(ls); err != nil {
		fmt.Println(fmt.Sprintln("Error While runing server: %v", err))
		log.Fatalf("failed to serve: %s", err)
	}

	fmt.Println("Server is running!")

}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "secret"
	dbName := "users"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:33060)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

type config struct {
	Port int `env:"GRPCSERVICE_PORT" envDefault:"9000"`
}
