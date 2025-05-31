package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	api "github.com/faisolarifin/wacoregateway/http/grpc"
	"github.com/faisolarifin/wacoregateway/model/constant"
	"github.com/faisolarifin/wacoregateway/provider"
	"github.com/faisolarifin/wacoregateway/provider/messaging"
	"github.com/faisolarifin/wacoregateway/service"
	"github.com/faisolarifin/wacoregateway/util"
	"github.com/go-playground/validator/v10"
)

func init() {
	if err := util.LoadConfig("."); err != nil {
		log.Fatal(err)
	}
}

func main() {
	logger := provider.NewLogger()
	validate := validator.New()
	container, err := provider.SqlStoreContainer()
	if err != nil {
		logger.Errorf(provider.AppLog, "Failed to connect to database:", err)
	}
	conn, err := provider.NewAMQPConn()
	if err != nil {
		log.Fatal(err)
	}
	publisher := messaging.NewAMQPPublisher(conn)
	logger.Infof(provider.AppLog, "Application started")

	ctx := context.WithValue(context.Background(), constant.CtxReqIDKey, "MAIN")
	go func(logger provider.ILogger) {
		service := service.NewService(container, logger, publisher)
		if err := service.LoadClients(ctx, container); err != nil {
			logger.Errorf(provider.AppLog, "Failed to load new client: %v", err)
		}

		app := api.NewApp(validate, logger, container, service)
		server, err := app.GRPCServer()
		if err != nil {
			log.Fatal(err)
		}

		addr := fmt.Sprintf(":%v", util.Configuration.Server.Port)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		logger.Infof(provider.AppLog, "gRPC server listening on %v", addr)
		if err := server.Serve(lis); err != nil {
			logger.Errorf(provider.AppLog, "failed to serve: %v", err)
		}
	}(logger)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownCh
	logger.Infof(provider.AppLog, "Receiving signal: %s", sig)

	func(logger provider.ILogger) {
		defer container.Close()

		logger.Infof(provider.AppLog, "Successfully stop Application.")
	}(logger)
}
