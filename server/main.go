package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

	var userGrpc *grpc.ClientConn
	var err error
	{
		userAdrr := fmt.Sprintf("%v:%v", cts.UserHost, cts.UserPort)
		userGrpc, err = grpc.Dial(userAdrr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			level.Error(logger).Log("gRPC: ", err)
			os.Exit(-1)
		}
	}

	opt := runtime.WithErrorHandler(errorHandler)
	gmux := runtime.NewServeMux(opt)
	ctx := context.Background()
	err = userpb.RegisterUserServiceHandler(ctx, gmux, userGrpc)

	if err != nil {
		level.Error(logger).Log("gRPC: ", err)
		os.Exit(-1)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		fmt.Println("Listenign on port: 8080")
		gwServer := &http.Server{
			Addr:    ":8080",
			Handler: gmux,
		}
		errs <- gwServer.ListenAndServe()
	}()

	level.Error(logger).Log("exit: ", <-errs)
}

func errorHandler(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	e, ok := status.FromError(err)
	if !ok {
		panic("Unsupported error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errors.ResolveHTTP(e.Code()))
	json.NewEncoder(w).Encode(map[string]interface{}{"error": e.Message()})
}

type constants struct {
	UserHost    string `env:"USER_SERVER,required"`
	UserPort    int    `env:"USER_PORT" envDefault:"50051"`
	DetailsHost string `env:"DETAILS_SERVER,required"`
	DetailsPort int    `env:"DETAILS_PORT" envDefault:"50051"`
}
