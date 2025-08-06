package service

import (
	"context"
	"fmt"

	"wacoregateway/internal/provider/messaging"
	"wacoregateway/util"

	"github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func AttachAllHandlers(deviceID string, publisher messaging.AMQPPublisherInterface, client *whatsmeow.Client) {
	client.AddEventHandler(func(evt interface{}) {
		HandleConnectionEvents(deviceID, evt)
		HandleMessageEvents(deviceID, publisher, evt)
		HandleQREvents(deviceID, evt)
		HandleAnyEvents(deviceID, evt)
	})
}

func HandleQREvents(deviceID string, evt interface{}) {
	switch v := evt.(type) {
	case *events.QR:
		logrus.Infof("[%s] QR: %s", deviceID, v.Codes)

	case *events.Receipt:
		fmt.Println("Receipt for message ID %s from %s ",
			v.MessageIDs, v.Sender.String())
	}
}

func HandleConnectionEvents(deviceID string, evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		logrus.Infof("[%s] Connected", deviceID)
	case *events.Disconnected:
		logrus.Warnf("[%s] Disconnected: %v", deviceID, v)
	case *events.LoggedOut:
		logrus.Warnf("[%s] Logged out", deviceID)
	}
}

func HandleMessageEvents(deviceID string, publisher messaging.AMQPPublisherInterface, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		sender := v.Info.Sender.String()
		content := v.Message.GetConversation()

		rawData := map[string]interface{}{
			"deviceID": deviceID,
			"sender":   sender,
			"content":  content,
			"metadata": map[string]interface{}{},
		}

		err := publisher.Publish(context.Background(), util.Configuration.AMQP.WaCoreGatewayQueue, rawData, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			fmt.Println("Failed to publish message: %v", err)
		}

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
	}
}

func HandleAnyEvents(deviceID string, evt interface{}) {
	switch v := evt.(type) {
	case *events.PairSuccess:
		logrus.Infof("[%s] Paired with %s", deviceID)

	case *events.MediaRetryError:
		fmt.Println("Ack for message %s - status: %s")

	case *events.Presence:
		fmt.Println("Presence %s:", v.From.String())

	case *events.CallOffer:
		fmt.Println("Call offer from")
	}
}
