package grpcHandlers

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
)

func (h *Handlers) CreateComment(ctx context.Context,
	request *pb.CreateCommentRequest) (*pb.GetCommentResponse, error) {
	commentID, err := h.service.CreateComment(ctx, models.CreateCommentPayload{
		Content:  request.GetContent(),
		AuthorID: request.GetCreatedBy(),
		PostID:   request.GetPost_ID(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create comment: %w", err)
	}

	comment, err := h.service.GetComment(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("cannot get created comment: %w", err)
	}

	return &pb.GetCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedBy: comment.CreatedBy,
		CreatedAt: comment.CreatedAt.String(),
		Post_ID:   comment.PostID,
	}, nil
}

func (h *Handlers) GetComment(ctx context.Context, request *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	comment, err := h.service.GetComment(ctx, request.GetID())
	if err != nil {
		return nil, fmt.Errorf("cannot get comment: %w", err)
	}

	return &pb.GetCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedBy: comment.CreatedBy,
		CreatedAt: comment.CreatedAt.String(),
		Post_ID:   comment.PostID,
	}, nil
}

func (h *Handlers) UpdateComment(ctx context.Context,
	request *pb.UpdateCommentRequest) (*pb.GetCommentResponse, error) {
	err := h.service.UpdateComment(ctx, models.UpdateCommentPayload{
		CommentID: request.GetID(),
		Content:   request.GetContent(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update comment: %w", err)
	}

	comment, err := h.service.GetComment(ctx, request.GetID())
	if err != nil {
		return nil, fmt.Errorf("cannot get updated comment: %w", err)
	}

	return &pb.GetCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedBy: comment.CreatedBy,
		CreatedAt: comment.CreatedAt.String(),
		Post_ID:   comment.PostID,
	}, nil
}

func (h *Handlers) DeleteComment(ctx context.Context,
	request *pb.GetCommentRequest) (*pb.DeleteCommentResponse, error) {
	if err := h.service.DeleteComment(ctx, request.GetID()); err != nil {
		return nil, fmt.Errorf("cannot delete comment: %w", err)
	}

	return &pb.DeleteCommentResponse{}, nil
}

func (h *Handlers) GetListOfComments(ctx context.Context,
	request *pb.GetListOfCommentsRequest) (*pb.GetListOfCommentsResponse, error) {
	comments, err := h.service.ListComments(ctx, models.Pagination{
		Limit:  uint64(request.GetPagination().GetLimit()),
		Offset: uint64(request.GetPagination().GetOffset()),
	}, models.FilterComments{
		Field: models.FilterCommentsByField(request.GetFilterComments().GetField()),
		Value: request.GetFilterComments().GetValue(),
	}, models.SortComments{
		Field: models.SortCommentsByField(request.GetSortComments().GetField()),
		IsASC: request.GetSortComments().GetIsAsc(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot get list of comments, %w", err)
	}

	var outs pb.GetListOfCommentsResponse

	for _, comment := range comments {
		out := pb.GetCommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedBy: comment.CreatedBy,
			CreatedAt: comment.CreatedAt.String(),
			Post_ID:   comment.PostID,
		}

		outs.Comments = append(outs.Comments, &out)
	}

	return &outs, nil
}
