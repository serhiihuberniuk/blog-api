package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/graphql/graph/generated"
	"github.com/serhiihuberniuk/blog-api/view/graphql/graph/model"
)

func (r *commentResolver) CreatedBy(ctx context.Context, obj *model.Comment) (*model.User, error) {
	user, err := r.Query().GetUser(ctx, obj.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("cannot get author of comment: %w", err)
	}

	return user, nil
}

func (r *commentResolver) Post(ctx context.Context, obj *model.Comment) (*model.Post, error) {
	post, err := r.Query().GetPost(ctx, obj.PostID)
	if err != nil {
		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	return post, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	userId, err := r.service.CreateUser(ctx, models.CreateUserPayload{
		Name:  input.Name,
		Email: input.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create user, %w", err)
	}

	user, err := r.Query().GetUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("cannot get created user:%w", err)
	}

	return user, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string,
	input model.UpdateUserInput) (*model.User, error) {
	err := r.service.UpdateUser(ctx, models.UpdateUserPayload{
		UserID: id,
		Name:   input.Name,
		Email:  input.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update user: %w", err)
	}

	user, err := r.Query().GetUser(ctx, id)
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

func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	postId, err := r.service.CreatePost(ctx, models.CreatePostPayload{
		Title:       input.Title,
		Description: input.Description,
		AuthorID:    input.CreatedBy,
		Tags:        input.Tags,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create post: %w", err)
	}

	post, err := r.Query().GetPost(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("cannot get created post: %w", err)
	}

	return post, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string,
	input model.UpdatePostInput) (*model.Post, error) {
	err := r.service.UpdatePost(ctx, models.UpdatePostPayload{
		PostID:      id,
		Title:       input.Title,
		Description: input.Description,
		Tags:        input.Tags,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update post: %w", err)
	}

	post, err := r.Query().GetPost(ctx, id)
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

func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error) {
	commentID, err := r.service.CreateComment(ctx, models.CreateCommentPayload{
		Content:  input.Content,
		PostID:   input.PostID,
		AuthorID: input.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create comment: %w", err)
	}

	comment, err := r.Query().GetComment(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("cannot get created comment: %w", err)
	}

	return comment, nil
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string,
	input model.UpdateCommentInput) (*model.Comment, error) {
	err := r.service.UpdateComment(ctx, models.UpdateCommentPayload{
		CommentID: id,
		Content:   input.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update comment: %w", err)
	}

	comment, err := r.Query().GetComment(ctx, id)
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

func (r *postResolver) CreatedBy(ctx context.Context, obj *model.Post) (*model.User, error) {
	user, err := r.Query().GetUser(ctx, obj.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("cannot get author of post: %w", err)
	}

	return user, nil
}

func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	user, err := r.service.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get user: %w", err)
	}

	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func (r *queryResolver) GetPost(ctx context.Context, id string) (*model.Post, error) {
	post, err := r.service.GetPost(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get post: %w", err)
	}

	out := &model.Post{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		AuthorID:    post.CreatedBy,
		CreatedAt:   post.CreatedAt.String(),
		Tags:        post.Tags,
	}

	out.CreatedBy, err = r.Post().CreatedBy(ctx, out)
	if err != nil {
		return nil, fmt.Errorf("cannot get author of post: %w", err)
	}

	return out, nil
}

func (r *queryResolver) GetComment(ctx context.Context, id string) (*model.Comment, error) {
	comment, err := r.service.GetComment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get comment, %w", err)
	}

	out := &model.Comment{
		ID:        comment.ID,
		Content:   comment.Content,
		AuthorID:  comment.CreatedBy,
		CreatedAt: comment.CreatedAt.String(),
		PostID:    comment.PostID,
	}

	out.CreatedBy, err = r.Comment().CreatedBy(ctx, out)
	if err != nil {
		return nil, fmt.Errorf("cannot get author of comment: %w", err)
	}

	out.Post, err = r.Comment().Post(ctx, out)
	if err != nil {
		return nil, fmt.Errorf("cannot get parent post: %w", err)
	}

	return out, nil
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
			sortPost.IsASC = sortPostsInput.IsAsc
		}
	}

	posts, err := r.service.ListPosts(ctx, pagination, filterPost, sortPost)
	if err != nil {
		return nil, fmt.Errorf("cannot get list of posts: %w", err)
	}

	listPosts := make([]*model.Post, 0, len(posts))

	for _, post := range posts {
		out, err := r.Query().GetPost(ctx, post.ID)
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
			sortComments.IsASC = sortCommentsInput.IsAsc
		}
	}

	comments, err := r.service.ListComments(ctx, pagination, filterComments, sortComments)
	if err != nil {
		return nil, fmt.Errorf("cannot get list of comments: %w", err)
	}

	listComments := make([]*model.Comment, 0, len(comments))

	for _, comment := range comments {
		out, err := r.Query().GetComment(ctx, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("cannot get comment: %w", err)
		}

		listComments = append(listComments, out)
	}

	return listComments, nil
}

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Post returns generated.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type (
	commentResolver  struct{ *Resolver }
	mutationResolver struct{ *Resolver }
	postResolver     struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)
