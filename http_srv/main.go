package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/go-kit/log"
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

	var userGRPC, detailsGRPC *grpc.ClientConn
	var grpcErr error
	{
		// userGRPC
		userAddr := fmt.Sprintf("%v:%v", cts.UserHost, cts.UserPort)
		userGRPC, grpcErr = grpc.Dial(userAddr, grpc.WithInsecure())
		if grpcErr != nil {
			level.Error(logger).Log("gRPC", grpcErr)
			os.Exit(-1)
		}

		// detailsGRPC
		detailsAddr := fmt.Sprintf("%v:%v", cts.DetailsHost, cts.DetailsPort)
		detailsGRPC, grpcErr = grpc.Dial(detailsAddr, grpc.WithInsecure())
		if grpcErr != nil {
			level.Error(logger).Log("gRPC", grpcErr)
			os.Exit(-1)
		}
	}

	ctx := context.Background()
	var httpSrv service.HTTPServicer
	{
		repository := repository.NewHTTPRepository(userGRPC, detailsGRPC, logger)
		httpSrv = service.NewHTTPService(repository, logger)
	}

	err := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		err <- fmt.Errorf("%s", <-c)
	}()

	httpEndpoints := transport.MakeHTTPEndpoints(httpSrv)

	go func() {
		fmt.Println("Listengin on port: 8080")
		httpHandler := transport.NewHTTPServer(ctx, httpEndpoints)
		err <- http.ListenAndServe(":8080", httpHandler)
	}()

	level.Error(logger).Log("exit: ", <-err)
}

type constants struct {
	UserHost    string `env:"USER_SERVER,required"`
	UserPort    int    `env:"USER_PORT" envDefault:"50051"`
	DetailsHost string `env:"DETAILS_SERVER,required"`
	DetailsPort int    `env:"DETAILS_PORT" envDefault:"50051"`
}
