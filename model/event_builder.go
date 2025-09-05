package model

import (
	"time"

	"github.com/google/uuid"
	"go.mau.fi/whatsmeow/types/events"
)

// EventBuilder helps create QueueEvent instances for different event types
type EventBuilder struct {
	SenderJID string
}

// NewEventBuilder creates a new EventBuilder instance
func NewEventBuilder(SenderJID string) *EventBuilder {
	return &EventBuilder{
		SenderJID: SenderJID,
	}
}

// CreateConnectedEvent creates a queue event for connected events
func (eb *EventBuilder) CreateConnectedEvent() *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeConnected,
		Timestamp: time.Now(),
		Data: ConnectionEventData{
			Status: "connected",
		},
	}
}

// CreateDisconnectedEvent creates a queue event for disconnected events
func (eb *EventBuilder) CreateDisconnectedEvent(reason string) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeDisconnected,
		Timestamp: time.Now(),
		Data: ConnectionEventData{
			Status: "disconnected",
			Reason: reason,
		},
	}
}

// CreateLoggedOutEvent creates a queue event for logged out events
func (eb *EventBuilder) CreateLoggedOutEvent() *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeLoggedOut,
		Timestamp: time.Now(),
		Data: ConnectionEventData{
			Status: "logged_out",
		},
	}
}

// CreateQREvent creates a queue event for QR events
func (eb *EventBuilder) CreateQREvent(code string) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeQR,
		Timestamp: time.Now(),
		Data: QREventData{
			Code: code,
		},
	}
}

// CreateTextMessageEvent creates a queue event for text message events
func (eb *EventBuilder) CreateTextMessageEvent(sender, content string) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeInboundMessage,
		Timestamp: time.Now(),
		Data: MessageEventData{
			Sender:      sender,
			MessageType: MessageTypeText,
			Content:     content,
			Metadata:    map[string]interface{}{},
		},
	}
}

// CreateImageMessageEvent creates a queue event for image message events
func (eb *EventBuilder) CreateImageMessageEvent(sender, caption, mimeType string, fileSize uint64) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeInboundMessage,
		Timestamp: time.Now(),
		Data: ImageMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeImage,
				Metadata:    map[string]interface{}{},
			},
			Caption:  caption,
			MimeType: mimeType,
			FileSize: fileSize,
		},
	}
}

// CreateAudioMessageEvent creates a queue event for audio message events
func (eb *EventBuilder) CreateAudioMessageEvent(sender string, duration uint32, mimeType string, fileSize uint64) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeInboundMessage,
		Timestamp: time.Now(),
		Data: AudioMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeAudio,
				Metadata:    map[string]interface{}{},
			},
			Duration: duration,
			MimeType: mimeType,
			FileSize: fileSize,
		},
	}
}

// CreateVideoMessageEvent creates a queue event for video message events
func (eb *EventBuilder) CreateVideoMessageEvent(sender, caption, mimeType string, duration uint32, fileSize uint64) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeInboundMessage,
		Timestamp: time.Now(),
		Data: VideoMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeVideo,
				Metadata:    map[string]interface{}{},
			},
			Caption:  caption,
			Duration: duration,
			MimeType: mimeType,
			FileSize: fileSize,
		},
	}
}

// CreateDocumentMessageEvent creates a queue event for document message events
func (eb *EventBuilder) CreateDocumentMessageEvent(sender, fileName, mimeType string, fileSize uint64) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeInboundMessage,
		Timestamp: time.Now(),
		Data: DocumentMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeDocument,
				Metadata:    map[string]interface{}{},
			},
			FileName: fileName,
			MimeType: mimeType,
			FileSize: fileSize,
		},
	}
}

// CreateLocationMessageEvent creates a queue event for location message events
func (eb *EventBuilder) CreateLocationMessageEvent(sender string, latitude, longitude float64, name, address string) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeInboundMessage,
		Timestamp: time.Now(),
		Data: LocationMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeLocation,
				Metadata:    map[string]interface{}{},
			},
			Latitude:  latitude,
			Longitude: longitude,
			Name:      name,
			Address:   address,
		},
	}
}

// CreateReactionMessageEvent creates a queue event for reaction message events
func (eb *EventBuilder) CreateReactionMessageEvent(sender, text, targetKey, targetSender string) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeInboundMessage,
		Timestamp: time.Now(),
		Data: ReactionMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeReaction,
				Metadata:    map[string]interface{}{},
			},
			Text:         text,
			TargetKey:    targetKey,
			TargetSender: targetSender,
		},
	}
}

// CreateButtonResponseMessageEvent creates a queue event for button response message events
func (eb *EventBuilder) CreateButtonResponseMessageEvent(sender, selectedButtonID, displayText string) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeInboundMessage,
		Timestamp: time.Now(),
		Data: ButtonResponseMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeButton,
				Metadata:    map[string]interface{}{},
			},
			SelectedButtonID: selectedButtonID,
			DisplayText:      displayText,
		},
	}
}

// CreateListResponseMessageEvent creates a queue event for list response message events
func (eb *EventBuilder) CreateListResponseMessageEvent(sender, title, description, selectedRowID string) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeInboundMessage,
		Timestamp: time.Now(),
		Data: ListResponseMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeList,
				Metadata:    map[string]interface{}{},
			},
			Title:         title,
			Description:   description,
			SelectedRowID: selectedRowID,
		},
	}
}

// CreateReceiptEvent creates a queue event for receipt events
func (eb *EventBuilder) CreateReceiptEvent(messageIDs []string, sender string, timestamp int64) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeReceipt,
		Timestamp: time.Now(),
		Data: ReceiptEventData{
			MessageIDs: messageIDs,
			Sender:     sender,
			Timestamp:  timestamp,
		},
	}
}

// CreatePresenceEvent creates a queue event for presence events
func (eb *EventBuilder) CreatePresenceEvent(from, status string, timestamp int64) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypePresence,
		Timestamp: time.Now(),
		Data: PresenceEventData{
			From:      from,
			Status:    status,
			Timestamp: timestamp,
		},
	}
}

// CreateCallOfferEvent creates a queue event for call offer events
func (eb *EventBuilder) CreateCallOfferEvent(from, callID string, timestamp int64) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeCallOffer,
		Timestamp: time.Now(),
		Data: CallOfferEventData{
			From:      from,
			CallID:    callID,
			Timestamp: timestamp,
		},
	}
}

// CreateMediaRetryErrorEvent creates a queue event for media retry error events
func (eb *EventBuilder) CreateMediaRetryErrorEvent(messageID, errorMsg string) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeMediaRetryError,
		Timestamp: time.Now(),
		Data: MediaRetryErrorEventData{
			MessageID: messageID,
			Error:     errorMsg,
		},
	}
}

// CreatePairSuccessEvent creates a queue event for pair success events
func (eb *EventBuilder) CreatePairSuccessEvent(senderJID, PhoneNumber string, deviceInfo interface{}) *QueueEvent {
	// Put client cache
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: senderJID,
		EventType: EventTypePairSuccess,
		Timestamp: time.Now(),
		Data: PairSuccessEventData{
			AccountJID:  eb.SenderJID,
			DeviceInfo:  deviceInfo,
			PhoneNumber: PhoneNumber,
		},
	}
}

// CreateGenericMessageEvent creates a generic message event from events.Message
func (eb *EventBuilder) CreateGenericMessageEvent(evt *events.Message) *QueueEvent {
	sender := evt.Info.Sender.String()
	msg := evt.Message

	var messageData interface{}

	switch {
	case msg.GetConversation() != "":
		messageData = MessageEventData{
			Sender:      sender,
			MessageType: MessageTypeText,
			Content:     msg.GetConversation(),
			Metadata: map[string]interface{}{
				"message_id": evt.Info.ID,
				"timestamp":  evt.Info.Timestamp.Unix(),
				"from_me":    evt.Info.IsFromMe,
				"chat":       evt.Info.Chat.String(),
			},
		}

	case msg.GetImageMessage() != nil:
		imgMsg := msg.GetImageMessage()
		messageData = ImageMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeImage,
				Metadata: map[string]interface{}{
					"message_id": evt.Info.ID,
					"timestamp":  evt.Info.Timestamp.Unix(),
					"from_me":    evt.Info.IsFromMe,
					"chat":       evt.Info.Chat.String(),
				},
			},
			Caption:  imgMsg.GetCaption(),
			MimeType: imgMsg.GetMimetype(),
			FileSize: imgMsg.GetFileLength(),
			FileURL:  imgMsg.GetURL(),
		}

	case msg.GetAudioMessage() != nil:
		audioMsg := msg.GetAudioMessage()
		messageData = AudioMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeAudio,
				Metadata: map[string]interface{}{
					"message_id": evt.Info.ID,
					"timestamp":  evt.Info.Timestamp.Unix(),
					"from_me":    evt.Info.IsFromMe,
					"chat":       evt.Info.Chat.String(),
				},
			},
			Duration: audioMsg.GetSeconds(),
			MimeType: audioMsg.GetMimetype(),
			FileSize: audioMsg.GetFileLength(),
			FileURL:  audioMsg.GetURL(),
		}

	case msg.GetVideoMessage() != nil:
		videoMsg := msg.GetVideoMessage()
		messageData = VideoMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeVideo,
				Metadata: map[string]interface{}{
					"message_id": evt.Info.ID,
					"timestamp":  evt.Info.Timestamp.Unix(),
					"from_me":    evt.Info.IsFromMe,
					"chat":       evt.Info.Chat.String(),
				},
			},
			Caption:  videoMsg.GetCaption(),
			Duration: videoMsg.GetSeconds(),
			MimeType: videoMsg.GetMimetype(),
			FileSize: videoMsg.GetFileLength(),
			FileURL:  videoMsg.GetURL(),
		}

	case msg.GetDocumentMessage() != nil:
		docMsg := msg.GetDocumentMessage()
		messageData = DocumentMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeDocument,
				Metadata: map[string]interface{}{
					"message_id": evt.Info.ID,
					"timestamp":  evt.Info.Timestamp.Unix(),
					"from_me":    evt.Info.IsFromMe,
					"chat":       evt.Info.Chat.String(),
				},
			},
			FileName: docMsg.GetFileName(),
			MimeType: docMsg.GetMimetype(),
			FileSize: docMsg.GetFileLength(),
			FileURL:  docMsg.GetURL(),
		}

	case msg.GetLocationMessage() != nil:
		locMsg := msg.GetLocationMessage()
		messageData = LocationMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeLocation,
				Metadata: map[string]interface{}{
					"message_id": evt.Info.ID,
					"timestamp":  evt.Info.Timestamp.Unix(),
					"from_me":    evt.Info.IsFromMe,
					"chat":       evt.Info.Chat.String(),
				},
			},
			Latitude:  locMsg.GetDegreesLatitude(),
			Longitude: locMsg.GetDegreesLongitude(),
			Name:      locMsg.GetName(),
			Address:   locMsg.GetAddress(),
		}

	case msg.GetReactionMessage() != nil:
		reactionMsg := msg.GetReactionMessage()
		messageData = ReactionMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeReaction,
				Metadata: map[string]interface{}{
					"message_id": evt.Info.ID,
					"timestamp":  evt.Info.Timestamp.Unix(),
					"from_me":    evt.Info.IsFromMe,
					"chat":       evt.Info.Chat.String(),
				},
			},
			Text:      reactionMsg.GetText(),
			TargetKey: reactionMsg.GetKey().GetID(),
		}

	case msg.GetButtonsResponseMessage() != nil:
		buttonMsg := msg.GetButtonsResponseMessage()
		messageData = ButtonResponseMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeButton,
				Metadata: map[string]interface{}{
					"message_id": evt.Info.ID,
					"timestamp":  evt.Info.Timestamp.Unix(),
					"from_me":    evt.Info.IsFromMe,
					"chat":       evt.Info.Chat.String(),
				},
			},
			SelectedButtonID: buttonMsg.GetSelectedButtonID(),
			DisplayText:      buttonMsg.GetSelectedDisplayText(),
		}

	case msg.GetListResponseMessage() != nil:
		listMsg := msg.GetListResponseMessage()
		messageData = ListResponseMessageData{
			MessageEventData: MessageEventData{
				Sender:      sender,
				MessageType: MessageTypeList,
				Metadata: map[string]interface{}{
					"message_id": evt.Info.ID,
					"timestamp":  evt.Info.Timestamp.Unix(),
					"from_me":    evt.Info.IsFromMe,
					"chat":       evt.Info.Chat.String(),
				},
			},
			Title:         listMsg.GetTitle(),
			Description:   listMsg.GetDescription(),
			SelectedRowID: listMsg.GetSingleSelectReply().GetSelectedRowID(),
		}

	default:
		messageData = MessageEventData{
			Sender:      sender,
			MessageType: MessageTypeText,
			Content:     "Unknown message type",
			Metadata: map[string]interface{}{
				"message_id": evt.Info.ID,
				"timestamp":  evt.Info.Timestamp.Unix(),
				"from_me":    evt.Info.IsFromMe,
				"chat":       evt.Info.Chat.String(),
			},
		}
	}

	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeInboundMessage,
		Timestamp: time.Now(),
		Data:      messageData,
	}
}

// CreateOutboundMessageEvent creates a queue event for outbound message events
func (eb *EventBuilder) CreateOutboundMessageEvent(messageID, messageType, to string, message interface{}) *QueueEvent {
	return &QueueEvent{
		EventID:   uuid.New().String(),
		SenderJID: eb.SenderJID,
		EventType: EventTypeOutboundMessage,
		Timestamp: time.Now(),
		Data: OutboundMessageData{
			MessageID:   messageID,
			MessageType: messageType,
			To:          to,
			Message:     message,
		},
	}
}
