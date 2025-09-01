package model

import (
	"time"
)

// EventType represents the type of WhatsApp event
type EventType string

const (
	// Connection Events
	EventTypeConnected    EventType = "connected"
	EventTypeDisconnected EventType = "disconnected"
	EventTypeLoggedOut    EventType = "logged_out"
	EventTypePairSuccess  EventType = "pair_success"

	// Message Events
	EventTypeMessage EventType = "message"

	// QR Events
	EventTypeQR EventType = "qr"

	// Receipt Events
	EventTypeReceipt EventType = "receipt"

	// Presence Events
	EventTypePresence EventType = "presence"

	// Call Events
	EventTypeCallOffer EventType = "call_offer"

	// Media Events
	EventTypeMediaRetryError EventType = "media_retry_error"
)

// MessageType represents the type of message content
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeVideo    MessageType = "video"
	MessageTypeDocument MessageType = "document"
	MessageTypeLocation MessageType = "location"
	MessageTypeReaction MessageType = "reaction"
	MessageTypeButton   MessageType = "button_response"
	MessageTypeList     MessageType = "list_response"
)

// QueueEvent is the main structure for all events sent to RabbitMQ
type QueueEvent struct {
	EventID   string      `json:"event_id"`
	SenderJID string      `json:"sender_jid"`
	EventType EventType   `json:"event_type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// ConnectionEventData represents connection-related events
type ConnectionEventData struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// QREventData represents QR code events
type QREventData struct {
	Code string `json:"code"`
}

// MessageEventData represents message events
type MessageEventData struct {
	Sender      string                 `json:"sender"`
	MessageType MessageType            `json:"message_type"`
	Content     string                 `json:"content,omitempty"`
	Caption     string                 `json:"caption,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ImageMessageData represents image message specific data
type ImageMessageData struct {
	MessageEventData
	Caption  string `json:"caption,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	FileSize uint64 `json:"file_size,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
}

// AudioMessageData represents audio message specific data
type AudioMessageData struct {
	MessageEventData
	Duration uint32 `json:"duration,omitempty"` // in seconds
	MimeType string `json:"mime_type,omitempty"`
	FileSize uint64 `json:"file_size,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
}

// VideoMessageData represents video message specific data
type VideoMessageData struct {
	MessageEventData
	Caption  string `json:"caption,omitempty"`
	Duration uint32 `json:"duration,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	FileSize uint64 `json:"file_size,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
}

// DocumentMessageData represents document message specific data
type DocumentMessageData struct {
	MessageEventData
	FileName string `json:"file_name,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	FileSize uint64 `json:"file_size,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
}

// LocationMessageData represents location message specific data
type LocationMessageData struct {
	MessageEventData
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// ReactionMessageData represents reaction message specific data
type ReactionMessageData struct {
	MessageEventData
	Text         string `json:"text"`
	TargetKey    string `json:"target_key,omitempty"`
	TargetSender string `json:"target_sender,omitempty"`
}

// ButtonResponseMessageData represents button response message specific data
type ButtonResponseMessageData struct {
	MessageEventData
	SelectedButtonID string `json:"selected_button_id"`
	DisplayText      string `json:"display_text,omitempty"`
}

// ListResponseMessageData represents list response message specific data
type ListResponseMessageData struct {
	MessageEventData
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	SelectedRowID string `json:"selected_row_id,omitempty"`
}

// ReceiptEventData represents message receipt events
type ReceiptEventData struct {
	MessageIDs []string `json:"message_ids"`
	Sender     string   `json:"sender"`
	Timestamp  int64    `json:"timestamp"`
}

// PresenceEventData represents presence events
type PresenceEventData struct {
	From      string `json:"from"`
	Status    string `json:"status"` // available, unavailable, composing, etc.
	Timestamp int64  `json:"timestamp"`
}

// CallOfferEventData represents call offer events
type CallOfferEventData struct {
	From      string `json:"from"`
	CallID    string `json:"call_id,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// MediaRetryErrorEventData represents media retry error events
type MediaRetryErrorEventData struct {
	MessageID string `json:"message_id,omitempty"`
	Error     string `json:"error,omitempty"`
}

// PairSuccessEventData represents successful pairing events
type PairSuccessEventData struct {
	AccountJID  string      `json:"account_jid,omitempty"`
	DeviceInfo  interface{} `json:"device_info,omitempty"`
	PhoneNumber string      `json:"phone_number,omitempty"`
}
