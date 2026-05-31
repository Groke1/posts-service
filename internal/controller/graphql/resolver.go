package graphql

import (
	"context"
	"posts-service/internal/entity"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

//go:generate mockgen -source=resolver.go -destination=mocks/resolver_mocks.go -package=mocks
type (
	postService interface {
		CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
		GetPost(ctx context.Context, id int64) (*entity.Post, error)
		GetPosts(ctx context.Context, authorID *int64, limit, offset int) ([]*entity.Post, error)
		SetPostCommentsEnabled(ctx context.Context, id, authorID int64, enabled bool) (*entity.Post, error)
	}

	commentService interface {
		CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
		GetComment(ctx context.Context, id int64) (*entity.Comment, error)
		GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int) ([]*entity.Comment, error)
	}

	commentSubscriber interface {
		Subscribe(ctx context.Context, postID int64) (<-chan *entity.Comment, func())
	}
)

type Resolver struct {
	postService    postService
	commentService commentService
	commentSub     commentSubscriber
}

func NewResolver(
	postService postService,
	commentService commentService,
	commentSub commentSubscriber) *Resolver {
	return &Resolver{
		postService:    postService,
		commentService: commentService,
		commentSub:     commentSub,
	}
}
