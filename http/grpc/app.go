package api

import (
	proto "github.com/faisolarifin/wacoregateway/model/pb"
	"github.com/faisolarifin/wacoregateway/provider"
	"github.com/faisolarifin/wacoregateway/service"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type App struct {
	validate *validator.Validate
	log      provider.ILogger
	service  service.ServiceInterface
}

type server struct {
	proto.UnimplementedWaCoreGatewayServer
	service service.ServiceInterface
}

func NewApp(validate *validator.Validate, log provider.ILogger, service service.ServiceInterface) *App {
	return &App{validate: validate, log: log, service: service}
}

func (a *App) GRPCServer() (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	proto.RegisterWaCoreGatewayServer(grpcServer, &server{
		service: a.service,
	})

	return grpcServer, nil
}
