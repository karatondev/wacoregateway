package service

import (
	"context"
	"fmt"

	"github.com/faisolarifin/wacoregateway/cache"
	proto "github.com/faisolarifin/wacoregateway/model/pb"
	"github.com/faisolarifin/wacoregateway/provider"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		client.AddEventHandler(eventHandler)
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
	client.AddEventHandler(eventHandler)

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

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		fmt.Println("Client connected")

	case *events.Disconnected:
		fmt.Println("Client disconnected")

	case *events.LoggedOut:
		fmt.Println("[%s] Logged out")

	case *events.Message:
		sender := v.Info.Sender.String()
		content := v.Message.GetConversation()
		fmt.Println("New message from %s: %s\n", sender, content)

		msg := v.Message

		switch {
		case msg.GetConversation() != "":
			fmt.Println("Text: %s", msg.GetConversation())

		case msg.GetImageMessage() != nil:
			fmt.Println("Image with caption: %s", msg.GetImageMessage().GetCaption())

		case msg.GetAudioMessage() != nil:
			fmt.Println("Voice note duration: %d", msg.GetAudioMessage().GetSeconds())

		case msg.GetDocumentMessage() != nil:
			fmt.Println("Document: %s", msg.GetDocumentMessage().GetFileName())

		case msg.GetVideoMessage() != nil:
			fmt.Println("Video with caption: %s", msg.GetVideoMessage().GetCaption())

		case msg.GetButtonsResponseMessage() != nil:
			fmt.Println("Button clicked: %s", msg.GetButtonsResponseMessage().SelectedButtonID)

		case msg.GetListResponseMessage() != nil:
			fmt.Println("List selected: %s", msg.GetListResponseMessage().GetTitle())

		case msg.GetLocationMessage() != nil:
			loc := msg.GetLocationMessage()
			fmt.Println("Location shared: %f,%f", loc.GetDegreesLatitude(), loc.GetDegreesLongitude())

		case msg.GetReactionMessage() != nil:
			fmt.Println("Reacted with: %s", msg.GetReactionMessage().GetText())
		}

	case *events.QR:
		fmt.Println("Scan QR: %s", v)

	case *events.Receipt:
		fmt.Println("Receipt for message ID %s from %s ",
			v.MessageIDs, v.Sender.String())

	case *events.MediaRetryError:
		fmt.Println("Ack for message %s - status: %s")

	case *events.Presence:
		fmt.Println("Presence %s:", v.From.String())

	case *events.CallOffer:
		fmt.Println("Call offer from")
	}
}
