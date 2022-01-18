package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"google.golang.org/grpc"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
}

func NewServer() *UserServer {
	return &UserServer{}
}

func main() {

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

	var user_grpc *grpc.ClientConn
	var err error
	{
		user_addr := "0.0.0.0:50051"
		user_grpc, err = grpc.Dial(user_addr, grpc.WithInsecure())
		if err != nil {
			level.Error(logger).Log("gRPC: ", err)
			os.Exit(-1)
		}
	}

	gmux := runtime.NewServeMux()
	ctx := context.Background()
	err = userpb.RegisterUserServiceHandler(ctx, gmux, user_grpc)

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
