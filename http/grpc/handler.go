package api

import (
	"context"

	proto "github.com/faisolarifin/wacoregateway/model/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) GetClientContact(ctx context.Context, req *proto.ClientdataRequest) (*proto.ContactListResponse, error) {

	senderJID := req.SenderJid
	if senderJID == "" {
		return nil, status.Errorf(codes.PermissionDenied, "senderJID param cannot be empty")
	}
	result, err := s.service.ProcessGetContact(ctx, senderJID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *server) GetClientGroup(ctx context.Context, req *proto.ClientdataRequest) (*proto.GroupListResponse, error) {

	senderJID := req.SenderJid
	if senderJID == "" {
		return nil, status.Errorf(codes.PermissionDenied, "senderJID param cannot be empty")
	}
	result, err := s.service.ProcessGetGroup(ctx, senderJID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *server) GetAllDevice(ctx context.Context, empty *emptypb.Empty) (*proto.DeviceListResponse, error) {
	result, err := s.service.ProcessGetDevices(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *server) SendMessage(ctx context.Context, req *proto.MessagePayload) (*proto.MessageResponse, error) {
	if req.SenderJid == "" {
		return nil, status.Errorf(codes.PermissionDenied, "senderJID param cannot be empty")
	}
	if req.To == "" {
		return nil, status.Errorf(codes.PermissionDenied, "to param cannot be empty")
	}

	result, err := s.service.ProcessSendMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *server) StreamConnectDevice(req *proto.ConnectDeviceRequest, stream proto.WaCoreGateway_StreamConnectDeviceServer) error {

	err := s.service.ConnectDevice(context.Background(), s.container, req, stream)
	if err != nil {
		return err
	}

	return nil
}
