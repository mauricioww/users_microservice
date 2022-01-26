package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"github.com/mauricioww/user_microsrv/user_details_srv/repository"
	"github.com/mauricioww/user_microsrv/user_details_srv/service"
	"github.com/mauricioww/user_microsrv/user_details_srv/transport"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			"GRPC_USER_DETAILS",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}
	level.Info(logger).Log("mesg", "service started")

	defer level.Info(logger).Log("msg", "service ended")

	var db *mongo.Database

	{
		mongoURI := fmt.Sprintf("mongodb://%v:%v", cts.DbHost, cts.DbPort)
		credentials := options.Credential{
			Username: cts.DbUser,
			Password: cts.DbPwd,
		}
		clientOpts := options.Client().ApplyURI(mongoURI).SetAuth(credentials)
		client, err := mongo.Connect(context.Background(), clientOpts)

		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}

		db = client.Database(cts.DbName)
	}

	var srv service.GrpcUserDetailsServicer
	{
		mongoRepository := repository.NewUserDetailsRepository(db, logger)
		srv = service.NewGrpcUserDetailsService(mongoRepository, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := transport.MakeGrpcEndpoints(srv)
	grpcServer := transport.NewGrpcUserDetailsServer(endpoints)
	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		logger.Log("Error listening: ", err)
		os.Exit(-1)
	}

	go func() {
		server := grpc.NewServer()
		detailspb.RegisterUserDetailsServiceServer(server, grpcServer)
		if err := server.Serve(listener); err != nil {
			logger.Log("Error serving", err)
		}
		level.Info(logger).Log("info", "grpc server started")
	}()

	level.Error(logger).Log("exit: ", <-errs)
}

type constants struct {
	DbUser string `env:"DB_USER,required"`
	DbPwd  string `env:"DB_PASSWORD,required"`
	DbHost string `env:"DB_HOST,required"`
	DbPort int    `env:"DB_PORT" envDefault:"27017"`
	DbName string `env:"DB_NAME" envDefault:"grpc_details"`
}
