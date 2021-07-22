package grpcHandlers

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handlers) CreatePost(ctx context.Context, request *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	postID, err := h.service.CreatePost(ctx, models.CreatePostPayload{
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		AuthorID:    request.GetCreatedBy(),
		Tags:        request.GetTags(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create post: %w", err)
	}

	post, err := h.service.GetPost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("cannot get created post: %w", err)
	}

	return &pb.CreatePostResponse{
		Id:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		CreatedBy:   post.CreatedBy,
		CreatedAt:   timestamppb.New(post.CreatedAt),
		Tags:        post.Tags,
	}, nil
}

func (h *Handlers) GetPost(ctx context.Context, request *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	post, err := h.service.GetPost(ctx, request.GetId())
	if err != nil {
		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	return &pb.GetPostResponse{
		Id:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		CreatedBy:   post.CreatedBy,
		CreatedAt:   timestamppb.New(post.CreatedAt),
		Tags:        post.Tags,
	}, nil
}

func (h *Handlers) UpdatePost(ctx context.Context, request *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	err := h.service.UpdatePost(ctx, models.UpdatePostPayload{
		PostID:      request.GetId(),
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		Tags:        request.GetTags(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update post: %w", err)
	}

	post, err := h.service.GetPost(ctx, request.GetId())
	if err != nil {
		return nil, fmt.Errorf("cannot get updated post: %w", err)
	}

	return &pb.UpdatePostResponse{
		Id:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		CreatedBy:   post.CreatedBy,
		CreatedAt:   timestamppb.New(post.CreatedAt),
		Tags:        post.Tags,
	}, nil
}

func (h *Handlers) DeletePost(ctx context.Context, request *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	if err := h.service.DeletePost(ctx, request.GetId()); err != nil {
		return nil, fmt.Errorf("cannot delete post: %w", err)
	}

	return &pb.DeletePostResponse{}, nil
}

func (h *Handlers) ListPosts(ctx context.Context,
	request *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	pagination := models.Pagination{}

	if request.GetPagination() != nil {
		limit := request.GetPagination().GetLimit()
		if limit <= 0 || limit > maxLimit {
			limit = maxLimit
		}

		offset := request.GetPagination().GetOffset()
		if offset < 0 {
			offset = 0
		}

		pagination = models.Pagination{
			Limit:  uint64(limit),
			Offset: uint64(offset),
		}
	}

	filter := models.FilterPosts{}
	if request.GetFilter() != nil {
		filter = models.FilterPosts{
			Field: models.FilterPostsByField(request.GetFilter().GetField().String()),
			Value: request.GetFilter().GetValue(),
		}
	}

	sort := models.SortPosts{}
	if request.GetSort() != nil {
		sort = models.SortPosts{
			SortByField: models.SortPostsByField(request.GetSort().GetField().String()),
			IsASC:       request.GetSort().GetIsAsc(),
		}
	}

	posts, err := h.service.ListPosts(ctx, pagination, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("cannot get list of posts: %w", err)
	}

	var outs pb.ListPostsResponse

	for _, post := range posts {
		out := pb.GetPostResponse{
			Id:          post.ID,
			Title:       post.Title,
			Description: post.Description,
			CreatedAt:   timestamppb.New(post.CreatedAt),
			CreatedBy:   post.CreatedBy,
			Tags:        post.Tags,
		}

		outs.Posts = append(outs.Posts, &out)
	}

	return &outs, nil
}
