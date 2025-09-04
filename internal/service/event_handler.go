package service

import (
	"context"

	"wacoregateway/internal/provider"
	"wacoregateway/internal/provider/messaging"
	"wacoregateway/model"
	"wacoregateway/util"

	"wacoregateway/model/constant"
	proto "wacoregateway/model/pb"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func AttachAllHandlers(senderJid string, publisher messaging.AMQPPublisherInterface, logger provider.ILogger, client *whatsmeow.Client, stream proto.WaCoreGateway_StreamConnectDeviceServer) {
	eventBuilder := model.NewEventBuilder(senderJid)
	ctx := context.WithValue(context.Background(), constant.CtxReqIDKey, senderJid)

	client.AddEventHandler(func(evt interface{}) {
		HandleConnectionEvents(senderJid, publisher, logger, eventBuilder, stream, ctx, evt)
		HandleMessageEvents(senderJid, publisher, logger, eventBuilder, ctx, evt)
		HandleQREvents(publisher, logger, eventBuilder, ctx, evt)
		HandleAnyEvents(senderJid, publisher, logger, eventBuilder, ctx, evt)
	})
}

func HandleQREvents(publisher messaging.AMQPPublisherInterface, logger provider.ILogger, eventBuilder *model.EventBuilder, ctx context.Context, evt interface{}) {
	switch v := evt.(type) {
	case *events.QR:
		logger.Infofctx(provider.AppLog, ctx, "QR: %v", v.Codes)

		// Create and publish QR event
		queueName := util.Configuration.Queues.QRHandlerQueue
		queueEvent := eventBuilder.CreateQREvent(v.Codes[0])
		err := publisher.Publish(ctx, queueName, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish QR event: %v", err)
		}
	}
}

func HandleConnectionEvents(senderJid string, publisher messaging.AMQPPublisherInterface, logger provider.ILogger, eventBuilder *model.EventBuilder, stream proto.WaCoreGateway_StreamConnectDeviceServer, ctx context.Context, evt interface{}) {
	queueName := util.Configuration.Queues.EventHandlerQueue
	queueEvent := &model.QueueEvent{}

	switch v := evt.(type) {
	case *events.PairSuccess:
		logger.Infofctx(provider.AppLog, ctx, "Paired with device")

		// Create and publish pair success event
		queueEvent = eventBuilder.CreatePairSuccessEvent(v.ID.String(), v.ID.User, v)
		go func() {
			streamData := &proto.EventResponse{
				Type: "event",
				Desc: "Pairing successfully",
			}
			if stream != nil {
				if err := stream.Send(streamData); err != nil {
					logger.Errorfctx(provider.AppLog, ctx, false, "Failed to send event: %v", err)
				}
			}
		}()

	case *events.Connected:
		logger.Infofctx(provider.AppLog, ctx, "Connected")

		queueEvent = eventBuilder.CreateConnectedEvent()

	case *events.Disconnected:
		logger.Infofctx(provider.AppLog, ctx, "Disconnected: %v", v)

		// Create and publish disconnected event
		queueEvent = eventBuilder.CreateDisconnectedEvent("disconnected")

	case *events.LoggedOut:
		logger.Infofctx(provider.AppLog, ctx, "Logged out")

		// Create and publish logged out event
		queueEvent = eventBuilder.CreateLoggedOutEvent()

	default:
		return
	}

	err := publisher.Publish(ctx, queueName, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
	if err != nil {
		logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish connection event: %v", err)
	}

}

func HandleMessageEvents(senderJid string, publisher messaging.AMQPPublisherInterface, logger provider.ILogger, eventBuilder *model.EventBuilder, ctx context.Context, evt interface{}) {
	queueName := util.Configuration.Queues.MessagesEventQueue
	queueEvent := &model.QueueEvent{}

	switch v := evt.(type) {
	case *events.Receipt:
		logger.Infofctx(provider.AppLog, ctx, "Receipt for message ID %v from %s", v.MessageIDs, v.Sender.String())

		// Create and publish receipt event
		queueEvent = eventBuilder.CreateReceiptEvent(v.MessageIDs, v.Sender.String(), v.Timestamp.Unix())

	case *events.Message:
		sender := v.Info.Sender.String()
		content := v.Message.GetConversation()

		// Use the generic message event creator which handles all message types
		queueEvent = eventBuilder.CreateGenericMessageEvent(v)

		logger.Infofctx(provider.AppLog, ctx, "New message from %s: %s", sender, content)

		msg := v.Message
		switch {
		case msg.GetConversation() != "":
			logger.Debugfctx(provider.AppLog, ctx, "Text: %s", msg.GetConversation())

		case msg.GetImageMessage() != nil:
			logger.Debugfctx(provider.AppLog, ctx, "Image with caption: %s", msg.GetImageMessage().GetCaption())

		case msg.GetAudioMessage() != nil:
			logger.Debugfctx(provider.AppLog, ctx, "Voice note duration: %d", msg.GetAudioMessage().GetSeconds())

		case msg.GetDocumentMessage() != nil:
			logger.Debugfctx(provider.AppLog, ctx, "Document: %s", msg.GetDocumentMessage().GetFileName())

		case msg.GetVideoMessage() != nil:
			logger.Debugfctx(provider.AppLog, ctx, "Video with caption: %s", msg.GetVideoMessage().GetCaption())

		case msg.GetButtonsResponseMessage() != nil:
			logger.Debugfctx(provider.AppLog, ctx, "Button clicked: %s", msg.GetButtonsResponseMessage().GetSelectedButtonID())

		case msg.GetListResponseMessage() != nil:
			logger.Debugfctx(provider.AppLog, ctx, "List selected: %s", msg.GetListResponseMessage().GetTitle())

		case msg.GetLocationMessage() != nil:
			loc := msg.GetLocationMessage()
			logger.Debugfctx(provider.AppLog, ctx, "Location shared: %f,%f", loc.GetDegreesLatitude(), loc.GetDegreesLongitude())

		case msg.GetReactionMessage() != nil:
			logger.Debugfctx(provider.AppLog, ctx, "Reacted with: %s", msg.GetReactionMessage().GetText())
		}

	default:
		return
	}

	err := publisher.Publish(ctx, queueName, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
	if err != nil {
		logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish message event: %v", err)
	}

}

func HandleAnyEvents(senderJid string, publisher messaging.AMQPPublisherInterface, logger provider.ILogger, eventBuilder *model.EventBuilder, ctx context.Context, evt interface{}) {
	queueName := util.Configuration.Queues.EventHandlerQueue
	queueEvent := &model.QueueEvent{}

	switch v := evt.(type) {
	case *events.MediaRetryError:
		logger.Infofctx(provider.AppLog, ctx, "Media retry error")

		// Create and publish media retry error event
		queueEvent = eventBuilder.CreateMediaRetryErrorEvent("", "media retry error")

	case *events.Presence:
		logger.Debugfctx(provider.AppLog, ctx, "Presence %s", v.From.String())

		// Create and publish presence event
		queueEvent = eventBuilder.CreatePresenceEvent(v.From.String(), "", v.LastSeen.Unix())

	case *events.CallOffer:
		logger.Infofctx(provider.AppLog, ctx, "Call offer received")

		// Create and publish call offer event
		queueEvent := eventBuilder.CreateCallOfferEvent(v.CallCreator.String(), v.CallID, 0)
		err := publisher.Publish(ctx, util.Configuration.Queues.EventHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish call offer event: %v", err)
		}
	default:
		return
	}

	err := publisher.Publish(ctx, queueName, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
	if err != nil {
		logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish any event: %v", err)
	}
}
