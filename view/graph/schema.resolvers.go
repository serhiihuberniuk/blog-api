package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/graph/generated"
	"github.com/serhiihuberniuk/blog-api/view/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input *model.CreateUserInput) (*model.User, error) {
	userId, err := r.service.CreateUser(ctx, models.CreateUserPayload{
		Name:  input.Name,
		Email: input.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create user, %w", err)
	}

	user, err := getUser(ctx, r.Resolver, userId)
	if err != nil {
		return nil, fmt.Errorf("cannot get created user:%w", err)
	}

	return user, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string,
	input *model.UpdateUserInput) (*model.User, error) {
	err := r.service.UpdateUser(ctx, models.UpdateUserPayload{
		UserID: id,
		Name:   input.Name,
		Email:  input.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update user: %w", err)
	}

	user, err := getUser(ctx, r.Resolver, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get updated user: %w", err)
	}

	return user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	if err := r.service.DeleteUser(ctx, id); err != nil {
		return false, fmt.Errorf("cannot dalete user: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input *model.CreatePostInput) (*model.Post, error) {
	postId, err := r.service.CreatePost(ctx, models.CreatePostPayload{
		Title:       input.Title,
		Description: input.Description,
		AuthorID:    input.CreatedBy,
		Tags:        input.Tags,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create post: %w", err)
	}

	post, err := getPost(ctx, r.Resolver, postId)
	if err != nil {
		return nil, fmt.Errorf("cannot get created post: %w", err)
	}

	return post, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string,
	input *model.UpdatePostInput) (*model.Post, error) {
	err := r.service.UpdatePost(ctx, models.UpdatePostPayload{
		PostID:      id,
		Title:       input.Title,
		Description: input.Description,
		Tags:        input.Tags,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update post: %w", err)
	}

	post, err := getPost(ctx, r.Resolver, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get updated post: %w", err)
	}

	return post, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	if err := r.service.DeletePost(ctx, id); err != nil {
		return false, fmt.Errorf("cannot dalete post: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input *model.CreateCommentInput) (*model.Comment, error) {
	commentID, err := r.service.CreateComment(ctx, models.CreateCommentPayload{
		Content:  input.Content,
		PostID:   input.PostID,
		AuthorID: input.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create comment: %w", err)
	}

	comment, err := getComment(ctx, r.Resolver, commentID)
	if err != nil {
		return nil, fmt.Errorf("cannot get created comment: %w", err)
	}

	return comment, nil
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string,
	input *model.UpdateCommentInput) (*model.Comment, error) {
	err := r.service.UpdateComment(ctx, models.UpdateCommentPayload{
		CommentID: id,
		Content:   input.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update comment: %w", err)
	}

	comment, err := getComment(ctx, r.Resolver, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get updated comment: %w", err)
	}

	return comment, nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (bool, error) {
	if err := r.service.DeleteComment(ctx, id); err != nil {
		return false, fmt.Errorf("cannot dalete comment: %w", err)
	}

	return true, nil
}

func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	user, err := getUser(ctx, r.Resolver, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	return user, nil
}

func (r *queryResolver) GetPost(ctx context.Context, id string) (*model.Post, error) {
	post, err := getPost(ctx, r.Resolver, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	return post, nil
}

func (r *queryResolver) GetComment(ctx context.Context, id string) (*model.Comment, error) {
	comment, err := getComment(ctx, r.Resolver, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get comment, %w", err)
	}

	return comment, nil
}

func (r *queryResolver) ListPosts(ctx context.Context, paginationInput *model.PaginationInput,
	filterPostsInput *model.FilterPostInput, sortPostsInput *model.SortPostsInput) ([]*model.Post, error) {
	pagination := getPaginationParams(paginationInput)

	filterPost := models.FilterPosts{}

	if filterPostsInput != nil {
		if filterPostsInput.Field.IsValid() {
			filterPost = models.FilterPosts{
				Field: allowedFilterPostsFields[filterPostsInput.Field],
				Value: filterPostsInput.Value,
			}
		}
	}

	sortPost := models.SortPosts{
		IsASC: true,
	}

	if sortPostsInput != nil {
		if sortPostsInput.Field.IsValid() {
			sortPost.SortByField = allowedSortPostsFields[sortPostsInput.Field]
			if sortPostsInput.IsAsc != nil {
				sortPost.IsASC = *sortPostsInput.IsAsc
			}
		}
	}

	posts, err := r.service.ListPosts(ctx, pagination, filterPost, sortPost)
	if err != nil {
		return nil, fmt.Errorf("cannot get list of posts: %w", err)
	}

	listPosts := make([]*model.Post, 0, len(posts))

	for _, post := range posts {
		out, err := getPost(ctx, r.Resolver, post.ID)
		if err != nil {
			return nil, fmt.Errorf("cannot get post, %w", err)
		}

		listPosts = append(listPosts, out)
	}

	return listPosts, nil
}

func (r *queryResolver) ListComments(ctx context.Context, paginationInput *model.PaginationInput,
	filterCommentsInput *model.FilterCommentsInput,
	sortCommentsInput *model.SortCommentsInput) ([]*model.Comment, error) {
	pagination := getPaginationParams(paginationInput)

	filterComments := models.FilterComments{}

	if filterCommentsInput != nil {
		if filterCommentsInput.Field.IsValid() {
			filterComments = models.FilterComments{
				Field: allowedFilterCommentsFields[filterCommentsInput.Field],
				Value: filterCommentsInput.Value,
			}
		}
	}

	sortComments := models.SortComments{
		IsASC: true,
	}

	if sortCommentsInput != nil {
		if sortCommentsInput.Field.IsValid() {
			sortComments.Field = allowedSortCommentsFields[sortCommentsInput.Field]

			if sortCommentsInput.IsAsc != nil {
				sortComments.IsASC = *sortCommentsInput.IsAsc
			}
		}
	}

	comments, err := r.service.ListComments(ctx, pagination, filterComments, sortComments)
	if err != nil {
		return nil, fmt.Errorf("cannot get list of comments: %w", err)
	}

	listComments := make([]*model.Comment, 0, len(comments))

	for _, comment := range comments {
		out, err := getComment(ctx, r.Resolver, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("cannot get comment: %w", err)
		}

		listComments = append(listComments, out)
	}

	return listComments, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)
