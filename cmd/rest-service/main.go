package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/go-kit/log"

	"github.com/neidersalgado/go-grpc-api/cmd/rest-service/user"
)

func main() {
	conf := Config{}

	if err := env.Parse(&conf); err != nil {
		fmt.Printf("%+v\n", err)
	}

	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	proxy := user.NewProxyRepository(logger)
	var handler http.Handler
	{
		handler = user.MakeHTTPHandler(*proxy, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	logger.Log("exit", <-errs)
}

type Config struct {
	Port         int    `env:"RESTSERVER_PORT" envDefault:"8080"`
	Env          string `env:"RESTSERVER_ENV" envDefault:"TEST"`
	Hosts        string `env:"RESTSERVER_HOSTS" envDefault:"127.0.0.1:"`
	WriteTimeout int    `env:"RESTSERVER_WRITETIMEOUT" envDefault:"15"`
	ReadTimeout  int    `env:"RESTSERVER_WRITETIMEOUT" envDefault:"15"`
}
