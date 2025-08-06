package service

import (
	"context"

	"wacoregateway/internal/provider"
	"wacoregateway/internal/provider/messaging"
	proto "wacoregateway/model/pb"

	"go.mau.fi/whatsmeow/store/sqlstore"
)

type ServiceInterface interface {
	LoadClients(ctx context.Context, container *sqlstore.Container) error
	ProcessGetDevices(ctx context.Context) (*proto.DeviceListResponse, error)
	ProcessGetContact(ctx context.Context, senderJID string) (*proto.ContactListResponse, error)
	ProcessGetGroup(ctx context.Context, senderJID string) (*proto.GroupListResponse, error)
	ProcessSendMessage(ctx context.Context, req *proto.MessagePayload) (*proto.MessageResponse, error)
	ConnectDevice(ctx context.Context, container *sqlstore.Container, req *proto.ConnectDeviceRequest, stream proto.WaCoreGateway_StreamConnectDeviceServer) error
}

type service struct {
	container *sqlstore.Container
	logger    provider.ILogger
	publisher messaging.AMQPPublisherInterface
}

func NewService(container *sqlstore.Container, logger provider.ILogger, publisher messaging.AMQPPublisherInterface) ServiceInterface {
	return &service{
		container: container,
		logger:    logger,
		publisher: publisher,
	}
}
