package api

import (
	"context"
	"fmt"

	"github.com/faisolarifin/wacoregateway/cache"
	proto "github.com/faisolarifin/wacoregateway/model/pb"
	"github.com/faisolarifin/wacoregateway/provider"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type App struct {
	validate *validator.Validate
	log      provider.ILogger
}

type server struct {
	proto.UnimplementedWaCoreGatewayServer
}

func NewApp(validate *validator.Validate, log provider.ILogger) *App {
	return &App{validate: validate, log: log}
}

func (a *App) GRPCServer() (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	proto.RegisterWaCoreGatewayServer(grpcServer, &server{})

	return grpcServer, nil
}

func (s *server) GetClientContact(ctx context.Context, req *proto.ContactRequest) (*proto.ContactResponse, error) {

	clients := cache.GetAllClients()
	devices := []string{}
	for id := range clients {
		devices = append(devices, id)
	}

	fmt.Println(devices)

	return &proto.ContactResponse{
		Code:    "200",
		Message: "success",
	}, nil
}
