package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"wacoregateway/internal/handler"
	"wacoregateway/internal/provider"
	"wacoregateway/internal/provider/messaging"
	"wacoregateway/internal/service"
	"wacoregateway/model/constant"
	"wacoregateway/util"

	"github.com/go-playground/validator/v10"
)

func Run(cfg *util.Config) {
	ctx := context.WithValue(context.Background(), constant.CtxReqIDKey, "MAIN")

	logger := provider.NewLogger()
	validate := validator.New()

	container, err := provider.SqlStoreContainer()
	if err != nil {
		logger.Errorfctx(provider.AppLog, ctx, false, "Failed to connect to database:", err)
	}
	conn, err := provider.NewAMQPConn()
	if err != nil {
		logger.Errorfctx(provider.AppLog, ctx, false, "Failed to connect to AMQP:", err)
	}
	publisher := messaging.NewAMQPPublisher(conn)
	logger.Infofctx(provider.AppLog, ctx, "Application started")

	go func(logger provider.ILogger) {
		service := service.NewService(container, logger, publisher)
		if err := service.LoadClients(ctx, container); err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to load new clients: %v", err)
		}

		app := handler.NewApp(validate, logger, container, service)
		server, err := app.GRPCServer()
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to create gRPC server: %v", err)
		}

		addr := fmt.Sprintf(":%v", util.Configuration.Server.Port)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "failed to listen: %v", err)
		}
		logger.Infofctx(provider.AppLog, ctx, "gRPC server listening on %v", addr)
		if err := server.Serve(lis); err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "failed to serve: %v", err)
		}
	}(logger)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownCh
	logger.Infofctx(provider.AppLog, ctx, "Receiving signal: %s", sig)

	func(logger provider.ILogger) {
		defer container.Close()

		logger.Infofctx(provider.AppLog, ctx, "Successfully stop Application.")
	}(logger)
}
