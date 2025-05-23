package service

import (
	"context"

	proto "github.com/faisolarifin/wacoregateway/model/pb"
	"github.com/faisolarifin/wacoregateway/provider"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type ServiceInterface interface {
	LoadNewClient(ctx context.Context, container *sqlstore.Container) error
	ProcessGetContact(ctx context.Context, senderJID string) (*proto.ContactListResponse, error)
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
