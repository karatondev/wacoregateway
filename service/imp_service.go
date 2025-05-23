package service

import (
	"context"
	"fmt"

	"github.com/faisolarifin/wacoregateway/cache"
	proto "github.com/faisolarifin/wacoregateway/model/pb"
	"github.com/faisolarifin/wacoregateway/provider"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) LoadNewClient(ctx context.Context, container *sqlstore.Container) error {
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
		result.Contacts = append(result.Contacts, &proto.ContactItem{
			Jid:   jid.String(),
			Name:  contact.FirstName,
			Short: contact.FullName,
		})
	}

	return result, nil
}
