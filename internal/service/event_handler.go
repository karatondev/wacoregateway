package service

import (
	"context"

	"wacoregateway/internal/provider"
	"wacoregateway/internal/provider/messaging"
	"wacoregateway/model"
	"wacoregateway/util"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func AttachAllHandlers(deviceID string, publisher messaging.AMQPPublisherInterface, logger provider.ILogger, client *whatsmeow.Client) {
	eventBuilder := model.NewEventBuilder(deviceID)
	ctx := context.Background()

	client.AddEventHandler(func(evt interface{}) {
		HandleConnectionEvents(deviceID, publisher, logger, eventBuilder, ctx, evt)
		HandleMessageEvents(deviceID, publisher, logger, eventBuilder, ctx, evt)
		HandleQREvents(deviceID, publisher, logger, eventBuilder, ctx, evt)
		HandleAnyEvents(deviceID, publisher, logger, eventBuilder, ctx, evt)
	})
}

func HandleQREvents(deviceID string, publisher messaging.AMQPPublisherInterface, logger provider.ILogger, eventBuilder *model.EventBuilder, ctx context.Context, evt interface{}) {
	switch v := evt.(type) {
	case *events.QR:
		logger.Infofctx(provider.AppLog, ctx, "[%s] QR: %v", deviceID, v.Codes)

		// Create and publish QR event
		queueEvent := eventBuilder.CreateQREvent(v.Codes[0])
		err := publisher.Publish(ctx, util.Configuration.Queues.QRHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish QR event: %v", err)
		}

	case *events.Receipt:
		logger.Infofctx(provider.AppLog, ctx, "Receipt for message ID %v from %s", v.MessageIDs, v.Sender.String())

		// Create and publish receipt event
		queueEvent := eventBuilder.CreateReceiptEvent(v.MessageIDs, v.Sender.String(), v.Timestamp.Unix())
		err := publisher.Publish(ctx, util.Configuration.Queues.EventHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish receipt event: %v", err)
		}
	}
}

func HandleConnectionEvents(deviceID string, publisher messaging.AMQPPublisherInterface, logger provider.ILogger, eventBuilder *model.EventBuilder, ctx context.Context, evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		logger.Infofctx(provider.AppLog, ctx, "[%s] Connected", deviceID)

		// Create and publish connected event
		queueEvent := eventBuilder.CreateConnectedEvent()
		err := publisher.Publish(ctx, util.Configuration.Queues.EventHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish connected event: %v", err)
		}

	case *events.Disconnected:
		logger.Infofctx(provider.AppLog, ctx, "[%s] Disconnected: %v", deviceID, v)

		// Create and publish disconnected event
		queueEvent := eventBuilder.CreateDisconnectedEvent("disconnected")
		err := publisher.Publish(ctx, util.Configuration.Queues.EventHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish disconnected event: %v", err)
		}

	case *events.LoggedOut:
		logger.Infofctx(provider.AppLog, ctx, "[%s] Logged out", deviceID)

		// Create and publish logged out event
		queueEvent := eventBuilder.CreateLoggedOutEvent()
		err := publisher.Publish(ctx, util.Configuration.Queues.EventHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish logged out event: %v", err)
		}
	}
}

func HandleMessageEvents(deviceID string, publisher messaging.AMQPPublisherInterface, logger provider.ILogger, eventBuilder *model.EventBuilder, ctx context.Context, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		sender := v.Info.Sender.String()
		content := v.Message.GetConversation()

		// Use the generic message event creator which handles all message types
		queueEvent := eventBuilder.CreateGenericMessageEvent(v)

		err := publisher.Publish(ctx, util.Configuration.Queues.EventHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish message event: %v", err)
		}

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
	}
}

func HandleAnyEvents(deviceID string, publisher messaging.AMQPPublisherInterface, logger provider.ILogger, eventBuilder *model.EventBuilder, ctx context.Context, evt interface{}) {
	switch v := evt.(type) {
	case *events.PairSuccess:
		logger.Infofctx(provider.AppLog, ctx, "[%s] Paired with device", deviceID)

		// Create and publish pair success event
		queueEvent := eventBuilder.CreatePairSuccessEvent("", "")
		err := publisher.Publish(ctx, util.Configuration.Queues.EventHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish pair success event: %v", err)
		}

	case *events.MediaRetryError:
		logger.Infofctx(provider.AppLog, ctx, "Media retry error")

		// Create and publish media retry error event
		queueEvent := eventBuilder.CreateMediaRetryErrorEvent("", "media retry error")
		err := publisher.Publish(ctx, util.Configuration.Queues.EventHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish media retry error event: %v", err)
		}

	case *events.Presence:
		logger.Debugfctx(provider.AppLog, ctx, "Presence %s", v.From.String())

		// Create and publish presence event
		queueEvent := eventBuilder.CreatePresenceEvent(v.From.String(), "", v.LastSeen.Unix())
		err := publisher.Publish(ctx, util.Configuration.Queues.EventHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish presence event: %v", err)
		}

	case *events.CallOffer:
		logger.Infofctx(provider.AppLog, ctx, "Call offer received")

		// Create and publish call offer event
		queueEvent := eventBuilder.CreateCallOfferEvent(v.CallCreator.String(), v.CallID, 0)
		err := publisher.Publish(ctx, util.Configuration.Queues.EventHandlerQueue, queueEvent, func(options *messaging.AMQPPublisherOptions) {})
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to publish call offer event: %v", err)
		}
	}
}
