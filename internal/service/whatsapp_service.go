package service

import (
	"context"

	"wacoregateway/internal/cache"
	"wacoregateway/internal/provider"
	proto "wacoregateway/model/pb"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) LoadClients(ctx context.Context, container *sqlstore.Container) error {
	clientLog := waLog.Stdout("Client", "DEBUG", true)

	devices, err := container.GetAllDevices(ctx)
	if err != nil {
		s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to get device store: %v", err)
		return err
	}

	for _, dev := range devices {
		client := whatsmeow.NewClient(dev, clientLog)
		err := client.Connect()
		if err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "failed to connect device %s: %v", dev.ID.String(), err)
			continue
		}
		AttachAllHandlers(dev.ID.String(), s.publisher, client)
		cache.SetClient(dev.ID.String(), client)
	}

	return nil
}

func (s *service) ConnectDevice(ctx context.Context, container *sqlstore.Container, req *proto.ConnectDeviceRequest, stream proto.WaCoreGateway_StreamConnectDeviceServer) error {
	clientLog := waLog.Stdout("Client", "DEBUG", true)

	user := req.Name
	jid := types.NewJID(user, types.DefaultUserServer)
	device := container.NewDevice()

	client := whatsmeow.NewClient(device, clientLog)
	AttachAllHandlers(device.ID.String(), s.publisher, client)

	cache.SetClient(jid.String(), client)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		go func() {
			_ = client.Connect()
		}()

		for evt := range qrChan {
			if evt.Event == whatsmeow.QRChannelEventCode {
				if err := stream.Send(&proto.EventResponse{
					Qr: evt.Code,
				}); err != nil {
					return err
				}
			}
		}
	} else {
		return status.Errorf(codes.AlreadyExists, "device with JID %s already exists", jid.String())
	}

	return nil
}

func (s *service) ProcessGetContact(ctx context.Context, senderJID string) (*proto.ContactListResponse, error) {
	client := cache.GetClient(senderJID)
	if client == nil {
		return nil, status.Errorf(codes.NotFound, "client with JID %s not found", senderJID)
	}

	contacts, err := client.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get contacts: %v", err)
	}

	result := &proto.ContactListResponse{}
	for jid, contact := range contacts {
		result.Contacts = append(result.Contacts, &proto.ClientdataItem{
			Jid:   jid.String(),
			Name:  contact.FirstName,
			Short: contact.FullName,
		})
	}

	return result, nil
}

func (s *service) ProcessGetGroup(ctx context.Context, senderJID string) (*proto.GroupListResponse, error) {
	client := cache.GetClient(senderJID)
	if client == nil {
		return nil, status.Errorf(codes.NotFound, "client with JID %s not found", senderJID)
	}

	groups, err := client.GetJoinedGroups()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get groups: %v", err)
	}

	result := &proto.GroupListResponse{}
	for _, group := range groups {
		result.Groups = append(result.Groups, &proto.ClientdataItem{
			Jid:  group.JID.String(),
			Name: group.Name,
		})
	}

	return result, nil
}

func (s *service) ProcessGetDevices(ctx context.Context) (*proto.DeviceListResponse, error) {

	clients := cache.GetAllClients()
	if clients == nil {
		return nil, status.Errorf(codes.NotFound, "no clients found")
	}
	result := &proto.DeviceListResponse{}
	for jid := range clients {
		result.Devices = append(result.Devices, &proto.DeviceItem{
			Jid: jid,
		})
	}

	return result, nil
}
