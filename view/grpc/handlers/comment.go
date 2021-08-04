package grpcHandlers

import (
	"context"

	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	allowedFilterCommentsFields = map[pb.ListCommentsRequest_Filter_Field]models.FilterCommentsByField{
		pb.ListCommentsRequest_Filter_UNKNOWN_FIELD: "",
		pb.ListCommentsRequest_Filter_POST_ID:       models.FilterCommentsByPost,
		pb.ListCommentsRequest_Filter_CREATED_AT:    models.FilterCommentsByCreatedAt,
		pb.ListCommentsRequest_Filter_CREATED_BY:    models.FilterCommentsByAuthor,
	}
	allowedSortCommentsFields = map[pb.ListCommentsRequest_Sort_Field]models.SortCommentsByField{
		pb.ListCommentsRequest_Sort_UNKNOWN_FIELD: "",
		pb.ListCommentsRequest_Sort_CREATED_AT:    models.SortCommentByCreatedAt,
	}
)

func (h *Handlers) CreateComment(ctx context.Context,
	request *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	commentID, err := h.service.CreateComment(ctx, models.CreateCommentPayload{
		Content:  request.GetContent(),
		AuthorID: request.GetCreatedBy(),
		PostID:   request.GetPostId(),
	})
	if err != nil {
		return nil, errorStatusGrpc(err)
	}

	comment, err := h.service.GetComment(ctx, commentID)
	if err != nil {
		return nil, errorStatusGrpc(err)
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
		return nil, errorStatusGrpc(err)
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
		return nil, errorStatusGrpc(err)
	}

	comment, err := h.service.GetComment(ctx, request.GetId())
	if err != nil {
		return nil, errorStatusGrpc(err)
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
		return nil, errorStatusGrpc(err)
	}

	return &pb.DeleteCommentResponse{}, nil
}

func (h *Handlers) ListComments(ctx context.Context,
	request *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	pagination := getPaginationParam(request.GetPagination())

	filter := models.FilterComments{}

	if request.GetFilter() != nil {
		filter = models.FilterComments{
			Field: allowedFilterCommentsFields[request.GetFilter().GetField()],
			Value: request.GetFilter().GetValue(),
		}
	}

	sort := models.SortComments{}

	if request.GetSort() != nil {
		sort = models.SortComments{
			Field: allowedSortCommentsFields[request.GetSort().GetField()],
			IsASC: request.GetSort().GetIsAsc(),
		}
	}

	comments, err := h.service.ListComments(ctx, pagination, filter, sort)
	if err != nil {
		return nil, errorStatusGrpc(err)
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
