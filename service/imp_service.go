package service

import (
	"context"
	"fmt"

	"github.com/faisolarifin/wacoregateway/cache"
	"github.com/faisolarifin/wacoregateway/provider"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
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
