package service

import (
	"context"
	"errors"
	"os"

	"github.com/faisolarifin/wacoregateway/cache"
	proto "github.com/faisolarifin/wacoregateway/model/pb"
	"go.mau.fi/whatsmeow"

	// waProto "go.mau.fi/whatsmeow/binary/proto"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// resp, err := client.SendMessage(context.Background(), jid, &waE2E.Message{
	// 	Conversation: gproto.String(req.MessageText),
	// })

	var msg *waProto.Message

	switch req.Type {

	case "text":
		msg = &waProto.Message{
			Conversation: protoStr(req.Text),
		}

	case "image":
		data, err := os.ReadFile(req.FilePath)
		if err != nil {
			return nil, err
		}
		uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaImage)
		if err != nil {
			return nil, err
		}
		msg = &waProto.Message{
			ImageMessage: &waProto.ImageMessage{
				URL:           &uploaded.URL,
				Mimetype:      protoStr("image/jpeg"),
				Caption:       protoStr(req.Text),
				FileSHA256:    uploaded.FileSHA256,
				FileEncSHA256: uploaded.FileEncSHA256,
				MediaKey:      uploaded.MediaKey,
				FileLength:    &uploaded.FileLength,
			},
		}

	case "video":
		data, err := os.ReadFile(req.FilePath)
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
				Mimetype:      protoStr("video/mp4"),
				Caption:       protoStr(req.Text),
				FileSHA256:    uploaded.FileSHA256,
				FileEncSHA256: uploaded.FileEncSHA256,
				MediaKey:      uploaded.MediaKey,
				FileLength:    &uploaded.FileLength,
			},
		}

	case "audio":
		data, err := os.ReadFile(req.FilePath)
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
				Mimetype:      protoStr("audio/ogg"),
				FileSHA256:    uploaded.FileSHA256,
				FileEncSHA256: uploaded.FileEncSHA256,
				MediaKey:      uploaded.MediaKey,
				FileLength:    &uploaded.FileLength,
			},
		}

	case "document":
		data, err := os.ReadFile(req.FilePath)
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
				Mimetype:      protoStr(req.MimeType),
				FileName:      protoStr(req.FileName),
				FileSHA256:    uploaded.FileSHA256,
				FileEncSHA256: uploaded.FileEncSHA256,
				MediaKey:      uploaded.MediaKey,
				FileLength:    &uploaded.FileLength,
			},
		}

	case "location":
		msg = &waProto.Message{
			LocationMessage: &waProto.LocationMessage{
				DegreesLatitude:  protoFloat(req.Latitude),
				DegreesLongitude: protoFloat(req.Longitude),
			},
		}

	case "button":
		msg = &waProto.Message{
			ButtonsMessage: &waProto.ButtonsMessage{
				ContentText: protoStr("Halo! Pilih tombol berikut:"),
				FooterText:  protoStr("Bot WhatsApp"),
				Buttons: []*waProto.ButtonsMessage_Button{
					{
						ButtonID:   protoStr("btn_1"),
						ButtonText: &waProto.ButtonsMessage_Button_ButtonText{DisplayText: protoStr("Tombol 1")},
						Type:       waProto.ButtonsMessage_Button_RESPONSE.Enum(),
					},
					{
						ButtonID:   protoStr("btn_2"),
						ButtonText: &waProto.ButtonsMessage_Button_ButtonText{DisplayText: protoStr("Tombol 2")},
						Type:       waProto.ButtonsMessage_Button_RESPONSE.Enum(),
					},
				},
				HeaderType: waProto.ButtonsMessage_UNKNOWN.Enum(),
			},
		}
	case "list":
		msg = &waProto.Message{
			ListMessage: &waProto.ListMessage{
				Title:       protoStr("Halo 👋"),
				Description: protoStr("Pilih salah satu menu di bawah ini:"),
				ButtonText:  protoStr("Lihat Menu"),
				Sections: []*waProto.ListMessage_Section{
					{
						Title: protoStr("Menu Utama"),
						Rows: []*waProto.ListMessage_Row{
							{
								RowID:       protoStr("menu_1"),
								Title:       protoStr("Layanan 1"),
								Description: protoStr("Deskripsi layanan 1"),
							},
							{
								RowID:       protoStr("menu_2"),
								Title:       protoStr("Layanan 2"),
								Description: protoStr("Deskripsi layanan 2"),
							},
						},
					},
				},
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
