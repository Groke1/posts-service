package post

import (
	"context"
	"posts-service/internal/entity"
)

//go:generate mockgen -source=post.go -destination=mocks/post_mocks.go -package=mocks

type (
	PostRepository interface {
		CreatePost(ctx context.Context, authorID int64, title string, text string, commentsEnabled bool) (*entity.Post, error)
		GetPostByID(ctx context.Context, id int64) (*entity.Post, error)
		GetPostsByAuthorID(ctx context.Context, authorID int64, limit, offset int) ([]*entity.Post, error)
		GetPosts(ctx context.Context, limit, offset int) ([]*entity.Post, error)
		SetPostCommentsEnabled(ctx context.Context, id int64, enabled bool) (*entity.Post, error)
	}
)

type postService struct {
	postRepository PostRepository
}

func NewPostService(postRepository PostRepository) *postService {
	return &postService{
		postRepository: postRepository,
	}
}

func (s *postService) CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	if post.Title == "" {
		return nil, entity.ErrPostTitleEmpty
	}
	if post.Text == "" {
		return nil, entity.ErrPostTextEmpty
	}
	return s.postRepository.CreatePost(ctx, post.AuthorID, post.Title, post.Text, post.CommentsEnabled)
}

func (s *postService) GetPost(ctx context.Context, id int64) (*entity.Post, error) {
	return s.postRepository.GetPostByID(ctx, id)
}

func (s *postService) GetPosts(ctx context.Context, authorID *int64, limit, offset int) ([]*entity.Post, error) {
	if limit <= 0 || offset < 0 {
		return nil, entity.ErrInvalidPagination
	}
	if authorID == nil {
		return s.postRepository.GetPosts(ctx, limit, offset)
	}
	return s.postRepository.GetPostsByAuthorID(ctx, *authorID, limit, offset)
}

func (s *postService) SetPostCommentsEnabled(ctx context.Context, id, authorID int64, enabled bool) (*entity.Post, error) {
	post, err := s.postRepository.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != authorID {
		return nil, entity.ErrForbidden
	}
	return s.postRepository.SetPostCommentsEnabled(ctx, id, enabled)
}
