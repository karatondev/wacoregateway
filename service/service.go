package service

import (
	"context"

	"github.com/faisolarifin/wacoregateway/provider"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type ServiceInterface interface {
	LoadNewClient(ctx context.Context, container *sqlstore.Container) error
}

type service struct {
	container *sqlstore.Container
	logger    provider.ILogger
}

func NewService(container *sqlstore.Container, logger provider.ILogger) ServiceInterface {
	return &service{
		container: container,
		logger:    logger,
	}
}
