package grpcHandlers

import (
	"context"
	"fmt"
	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handlers) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userID, err := h.service.CreateUser(ctx, models.CreateUserPayload{
		Name:  request.GetName(),
		Email: request.GetEmail(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create user: %w", err)
	}

	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get created user: %w", err)
	}

	return &pb.CreateUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *Handlers) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.service.GetUser(ctx, request.GetId())
	if err != nil {
		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	return &pb.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *Handlers) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	err := h.service.UpdateUser(ctx, models.UpdateUserPayload{
		UserID: request.GetId(),
		Name:   request.GetName(),
		Email:  request.GetEmail(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update user: %w", err)
	}

	user, err := h.service.GetUser(ctx, request.GetId())
	if err != nil {
		return nil, fmt.Errorf("cannot get updated user: %w", err)
	}

	return &pb.UpdateUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *Handlers) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := h.service.DeleteUser(ctx, request.GetId()); err != nil {
		return nil, fmt.Errorf("cannot delete user: %w", err)
	}

	return &pb.DeleteUserResponse{}, nil
}
