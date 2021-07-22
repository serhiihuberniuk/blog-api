package grpcHandlers

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
)

func (h *Handlers) CreateComment(ctx context.Context,
	request *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	commentID, err := h.service.CreateComment(ctx, models.CreateCommentPayload{
		Content:  request.GetContent(),
		AuthorID: request.GetCreatedBy(),
		PostID:   request.GetPostId(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create comment: %w", err)
	}

	comment, err := h.service.GetComment(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("cannot get created comment: %w", err)
	}

	return &pb.CreateCommentResponse{
		Id:        comment.ID,
		Content:   comment.Content,
		CreatedBy: comment.CreatedBy,
		CreatedAt: timestamppb.New(comment.CreatedAt),
		PostId:    comment.PostID,
	}, nil
}

func (h *Handlers) GetComment(ctx context.Context, request *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	comment, err := h.service.GetComment(ctx, request.GetId())
	if err != nil {
		return nil, fmt.Errorf("cannot get comment: %w", err)
	}

	return &pb.GetCommentResponse{
		Id:        comment.ID,
		Content:   comment.Content,
		CreatedBy: comment.CreatedBy,
		CreatedAt: timestamppb.New(comment.CreatedAt),
		PostId:    comment.PostID,
	}, nil
}

func (h *Handlers) UpdateComment(ctx context.Context,
	request *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	err := h.service.UpdateComment(ctx, models.UpdateCommentPayload{
		CommentID: request.GetId(),
		Content:   request.GetContent(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update comment: %w", err)
	}

	comment, err := h.service.GetComment(ctx, request.GetId())
	if err != nil {
		return nil, fmt.Errorf("cannot get updated comment: %w", err)
	}

	return &pb.UpdateCommentResponse{
		Id:        comment.ID,
		Content:   comment.Content,
		CreatedBy: comment.CreatedBy,
		CreatedAt: timestamppb.New(comment.CreatedAt),
		PostId:    comment.PostID,
	}, nil
}

func (h *Handlers) DeleteComment(ctx context.Context,
	request *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	if err := h.service.DeleteComment(ctx, request.GetId()); err != nil {
		return nil, fmt.Errorf("cannot delete comment: %w", err)
	}

	return &pb.DeleteCommentResponse{}, nil
}

func (h *Handlers) ListComments(ctx context.Context,
	request *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	comments, err := h.service.ListComments(ctx, models.Pagination{
		Limit:  uint64(request.GetPagination().GetLimit()),
		Offset: uint64(request.GetPagination().GetOffset()),
	}, models.FilterComments{
		Field: models.FilterCommentsByField(request.GetFilter().GetField().String()),
		Value: request.GetFilter().GetValue(),
	}, models.SortComments{
		Field: models.SortCommentsByField(request.GetSort().GetField().String()),
		IsASC: request.GetSort().GetIsAsc(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot get list of comments, %w", err)
	}

	var outs pb.ListCommentsResponse

	for _, comment := range comments {
		out := pb.GetCommentResponse{
			Id:        comment.ID,
			Content:   comment.Content,
			CreatedBy: comment.CreatedBy,
			CreatedAt: timestamppb.New(comment.CreatedAt),
			PostId:    comment.PostID,
		}

		outs.Comments = append(outs.Comments, &out)
	}

	return &outs, nil
}
