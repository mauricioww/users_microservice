package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
	"github.com/mauricioww/user_microsrv/http_srv/service"
	"github.com/mauricioww/user_microsrv/http_srv/transport"
	"google.golang.org/grpc"
)

func main() {
	cts := constants{}

	if err := env.Parse(&cts); err != nil {
		fmt.Printf("%+v\n", err)
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service",
			"HTTP_SRV",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}

	level.Info(logger).Log("mesg", "service started")

	defer level.Info(logger).Log("msg", "service ended")

	var user_grpc, details_grpc *grpc.ClientConn
	var grpc_err error
	{
		// user_grpc
		user_addr := fmt.Sprintf("%v:%v", cts.UserHost, cts.DetailsPort)
		user_grpc, grpc_err = grpc.Dial(user_addr, grpc.WithInsecure())
		if grpc_err != nil {
			level.Error(logger).Log("gRPC", grpc_err)
			os.Exit(-1)
		}

		// details_grpc
		details_addr := fmt.Sprintf("%v:%v", cts.DetailsHost, cts.DetailsPort)
		details_grpc, grpc_err = grpc.Dial(details_addr, grpc.WithInsecure())
		if grpc_err != nil {
			level.Error(logger).Log("gRPC", grpc_err)
			os.Exit(-1)
		}
	}

	ctx := context.Background()
	var http_srv service.HttpService
	{
		repository := repository.NewHttpRepository(user_grpc, details_grpc, logger)
		http_srv = service.NewHttpService(repository, logger)
	}

	err := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		err <- fmt.Errorf("%s", <-c)
	}()

	http_endpoints := transport.MakeHttpEndpoints(http_srv)

	go func() {
		fmt.Println("Listengin on port: 8080")
		http_handler := transport.NewHTTPServer(ctx, http_endpoints)
		err <- http.ListenAndServe(":8080", http_handler)
	}()

	level.Error(logger).Log("exit: ", <-err)
}

type constants struct {
	UserHost    string `env:"USER_SERVER,required"`
	UserPort    int    `env:"USER_PORT" envDefault:"50051"`
	DetailsHost string `env:"DETAILS_SERVER,required"`
	DetailsPort int    `env:"DETAILS_PORT" envDefault:"50051"`
}
