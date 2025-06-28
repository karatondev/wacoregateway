package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/faisolarifin/wacoregateway/cache"
	proto "github.com/faisolarifin/wacoregateway/model/pb"
	"go.mau.fi/whatsmeow"

	// waProto "go.mau.fi/whatsmeow/binary/proto"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	Text     = "text"
	Image    = "image"
	Video    = "video"
	Document = "document"
	Audio    = "audio"
	Location = "location"
)

func (s *service) ProcessSendMessage(ctx context.Context, req *proto.MessagePayload) (*proto.MessageResponse, error) {

	client := cache.GetClient(req.SenderJid)
	if client == nil {
		return nil, status.Errorf(codes.NotFound, "sender device with JID %s not found", req.SenderJid)
	}

	jid, err := types.ParseJID(req.To)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid recipient JID: %v", err)
	}

	var msg *waProto.Message

	switch req.Type {

	case Text:
		msg = &waProto.Message{
			Conversation: protoStr(req.Text),
		}

	case Image:
		url := req.Image.Url
		var data []byte

		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
			resp, err := http.Get(url)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to get image %v", err)
			}
			defer resp.Body.Close()

			data, err = io.ReadAll(resp.Body)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed get bytes the image: %v", err)
			}
		} else {
			data, err = os.ReadFile(req.Image.Url)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed get local image %v", err)
			}
		}

		uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaImage)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed upload image to whatsapp: %v", err)
		}
		msg = &waProto.Message{
			ImageMessage: &waProto.ImageMessage{
				URL:           &uploaded.URL,
				Mimetype:      protoStr(req.Image.Mimetype),
				Caption:       protoStr(req.Image.Caption),
				FileSHA256:    uploaded.FileSHA256,
				FileEncSHA256: uploaded.FileEncSHA256,
				MediaKey:      uploaded.MediaKey,
				FileLength:    &uploaded.FileLength,
				DirectPath:    protoStr(uploaded.DirectPath),
			},
		}

	case Video:
		data, err := os.ReadFile(req.Video.Url)
		if err != nil {
			return nil, err
		}
		uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaVideo)
		if err != nil {
			return nil, err
		}
		msg = &waProto.Message{
			VideoMessage: &waProto.VideoMessage{
				URL:           &uploaded.URL,
				Mimetype:      protoStr(req.Video.Mimetype),
				Caption:       protoStr(req.Video.Caption),
				FileSHA256:    uploaded.FileSHA256,
				FileEncSHA256: uploaded.FileEncSHA256,
				MediaKey:      uploaded.MediaKey,
				FileLength:    &uploaded.FileLength,
				DirectPath:    protoStr(uploaded.DirectPath),
			},
		}

	case Audio:
		data, err := os.ReadFile(req.Audio.Url)
		if err != nil {
			return nil, err
		}
		uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaAudio)
		if err != nil {
			return nil, err
		}
		msg = &waProto.Message{
			AudioMessage: &waProto.AudioMessage{
				URL:           &uploaded.URL,
				Mimetype:      protoStr(req.Audio.MimeType),
				FileSHA256:    uploaded.FileSHA256,
				FileEncSHA256: uploaded.FileEncSHA256,
				MediaKey:      uploaded.MediaKey,
				FileLength:    &uploaded.FileLength,
				PTT:           protoBool(req.Audio.Ptt),
				DirectPath:    protoStr(uploaded.DirectPath),
			},
		}

	case Document:
		data, err := os.ReadFile(req.Document.Url)
		if err != nil {
			return nil, err
		}
		uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaDocument)
		if err != nil {
			return nil, err
		}
		msg = &waProto.Message{
			DocumentMessage: &waProto.DocumentMessage{
				URL:           &uploaded.URL,
				Mimetype:      protoStr(req.Document.Mimetype),
				FileName:      protoStr(req.Document.Filename),
				FileSHA256:    uploaded.FileSHA256,
				FileEncSHA256: uploaded.FileEncSHA256,
				MediaKey:      uploaded.MediaKey,
				FileLength:    &uploaded.FileLength,
				Title:         protoStr(req.Document.Title),
				DirectPath:    protoStr(uploaded.DirectPath),
			},
		}

	case Location:
		msg = &waProto.Message{
			LocationMessage: &waProto.LocationMessage{
				DegreesLatitude:  protoFloat(req.Location.Latitude),
				DegreesLongitude: protoFloat(req.Location.Longitude),
				Name:             protoStr(req.Location.Name),
				Address:          protoStr(req.Location.Address),
			},
		}

	default:
		return nil, errors.New("unsupported message type")
	}

	resp, err := client.SendMessage(context.Background(), jid, msg)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	return &proto.MessageResponse{
		Id: resp.ID,
	}, nil
}

func protoStr(s string) *string {
	return &s
}

func protoBool(b bool) *bool {
	return &b
}

func protoFloat(f float64) *float64 {
	return &f
}
