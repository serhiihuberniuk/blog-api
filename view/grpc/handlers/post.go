package grpcHandlers

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/grpc/pb"
)

func (h *Handlers) CreatePost(ctx context.Context, request *pb.CreatePostRequest) (*pb.GetPostResponse, error) {
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

	return &pb.GetPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		CreatedBy:   post.CreatedBy,
		CreatedAt:   post.CreatedAt.String(),
		Tags:        post.Tags,
	}, nil
}

func (h *Handlers) GetPost(ctx context.Context, request *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	post, err := h.service.GetPost(ctx, request.GetID())
	if err != nil {
		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	return &pb.GetPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		CreatedBy:   post.CreatedBy,
		CreatedAt:   post.CreatedAt.String(),
		Tags:        post.Tags,
	}, nil
}

func (h *Handlers) UpdatePost(ctx context.Context, request *pb.UpdatePostRequest) (*pb.GetPostResponse, error) {
	err := h.service.UpdatePost(ctx, models.UpdatePostPayload{
		PostID:      request.GetID(),
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		Tags:        request.GetTags(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update post: %w", err)
	}

	post, err := h.service.GetPost(ctx, request.GetID())
	if err != nil {
		return nil, fmt.Errorf("cannot get updated post: %w", err)
	}

	return &pb.GetPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		CreatedBy:   post.CreatedBy,
		CreatedAt:   post.CreatedAt.String(),
		Tags:        post.Tags,
	}, nil
}

func (h *Handlers) DeletePost(ctx context.Context, request *pb.GetPostRequest) (*pb.DeletePostResponse, error) {
	if err := h.service.DeletePost(ctx, request.GetID()); err != nil {
		return nil, fmt.Errorf("cannot delete post: %w", err)
	}

	return &pb.DeletePostResponse{}, nil
}

func (h *Handlers) GetListOfPosts(ctx context.Context,
	request *pb.GetListOfPostsRequest) (*pb.GetListOfPostsResponse, error) {
	posts, err := h.service.ListPosts(ctx, models.Pagination{
		Limit:  uint64(request.GetPagination().GetLimit()),
		Offset: uint64(request.GetPagination().GetOffset()),
	},
		models.FilterPosts{
			Field: models.FilterPostsByField(request.GetFilterPosts().GetField()),
			Value: request.GetFilterPosts().GetValue(),
		},
		models.SortPosts{
			SortByField: models.SortPostsByField(request.GetSortPosts().GetField()),
			IsASC:       request.GetSortPosts().GetIsAsc(),
		})
	if err != nil {
		return nil, fmt.Errorf("cannot get list of posts: %w", err)
	}

	var outs pb.GetListOfPostsResponse

	for _, post := range posts {
		out := pb.GetPostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Description: post.Description,
			CreatedAt:   post.CreatedAt.String(),
			CreatedBy:   post.CreatedBy,
			Tags:        post.Tags,
		}

		outs.Posts = append(outs.Posts, &out)
	}

	return &outs, nil
}
