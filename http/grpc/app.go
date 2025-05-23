package api

import (
	"github.com/faisolarifin/wacoregateway/provider"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type App struct {
	validate *validator.Validate
	log      provider.ILogger
}

func NewApp(validate *validator.Validate, log provider.ILogger) *App {
	return &App{validate: validate, log: log}
}

func (a *App) GRPCServer() (*grpc.Server, error) {
	grpcServer := grpc.NewServer()

	return grpcServer, nil
}
