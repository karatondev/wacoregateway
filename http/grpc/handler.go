package api

import (
	"context"

	proto "github.com/faisolarifin/wacoregateway/model/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetClientContact(ctx context.Context, req *proto.ContactRequest) (*proto.ContactListResponse, error) {

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
