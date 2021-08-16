package grpcHandlers

import (
	"context"

	viewmodels "github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
)

func (h *Handlers) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := h.service.Login(ctx, viewmodels.LoginPayload{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, errorStatusGrpc(err)
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}
