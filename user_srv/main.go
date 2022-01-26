package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mauricioww/user_microsrv/user_srv/repository"
	"github.com/mauricioww/user_microsrv/user_srv/service"
	"github.com/mauricioww/user_microsrv/user_srv/transport"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
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
			"GRPC_SRV",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}

	level.Info(logger).Log("mesg", "service started")

	defer level.Info(logger).Log("msg", "service ended")

	var db *sql.DB
	{
		var err error
		mysqlAddr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", cts.DbUser, cts.DbPwd, cts.DbHost, cts.DbPort, cts.DbName)
		db, err = sql.Open("mysql", mysqlAddr)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	var srv service.GrpcUserServicer
	{
		mysqlRepo := repository.NewUserRepository(db, logger)
		srv = service.NewGrpcUserService(mysqlRepo, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := transport.MakeGrpcEndpoints(srv)
	grpcServer := transport.NewGrpcUserServer(endpoints)
	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		logger.Log("Error listening: ", err)
		os.Exit(-1)
	}

	go func() {
		server := grpc.NewServer()
		userpb.RegisterUserServiceServer(server, grpcServer)
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
	DbPort int    `env:"DB_PORT" envDefault:"3306"`
	DbName string `env:"DB_NAME" envDefault:"grpc_user"`
}
