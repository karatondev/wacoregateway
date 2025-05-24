package service

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/events"
	"github.com/faisolarifin/wacoregateway/cache"
	proto "github.com/faisolarifin/wacoregateway/model/pb"
	"github.com/faisolarifin/wacoregateway/provider"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gproto "google.golang.org/protobuf/proto"
)

func (s *service) LoadClients(ctx context.Context, container *sqlstore.Container) error {
	clientLog := waLog.Stdout("Client", "DEBUG", true)

	devices, err := container.GetAllDevices(ctx)
	if err != nil {
		s.logger.Errorf(provider.AppLog, "Failed to get device store: %v", err)
		return err
	}

	for _, dev := range devices {
		client := whatsmeow.NewClient(dev, clientLog)
		err := client.Connect()
		if err != nil {
			fmt.Printf("failed to connect device %s: %v\n", dev.ID.String(), err)
			continue
		}
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
	cache.SetClient(jid.String(), client)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		go func() {
			_ = client.Connect()
		}()

		for evt := range qrChan {
			if evt.Event == "code" {
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

func (s *service) ProcessSendMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {

	client := cache.GetClient(req.SenderJid)
	if client == nil {
		return nil, status.Errorf(codes.NotFound, "sender device with JID %s not found", req.SenderJid)
	}

	jid, err := types.ParseJID(req.To)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid recipient JID: %v", err)
	}

	resp, err := client.SendMessage(context.Background(), jid, &waE2E.Message{
		Conversation: gproto.String(req.MessageText),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}
	return &proto.MessageResponse{
		Id: resp.ID,
	}, nil
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Handle incoming messages
		fmt.Println(v)
		// sender := v.Info.Sender.String()
		// content := v.Message.GetConversation()
		// fmt.Printf("New message from %s: %s\n", sender, content)
	}
}
