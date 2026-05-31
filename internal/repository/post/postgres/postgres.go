package postgres

import (
	"context"
	"errors"
	"posts-service/internal/entity"
	postssqlc "posts-service/internal/repository/post/postgres/sqlc"

	"github.com/jackc/pgx/v5"
)

type postRepository struct {
	queries *postssqlc.Queries
}

func NewPostRepository(qdb postssqlc.DBTX) *postRepository {
	return &postRepository{
		queries: postssqlc.New(qdb),
	}
}

func (r *postRepository) GetPostByID(ctx context.Context, id int64) (*entity.Post, error) {
	row, err := r.queries.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrPostNotFound
		}
		return nil, err
	}
	return &entity.Post{
		ID:              row.ID,
		AuthorID:        row.AuthorID,
		Text:            row.Text,
		Title:           row.Title,
		CommentsEnabled: row.CommentsEnabled,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}, nil
}

func (r *postRepository) GetPostByIDForUpdate(ctx context.Context, id int64) (*entity.Post, error) {
	row, err := r.queries.GetPostForUpdate(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrPostNotFound
		}
		return nil, err
	}
	return &entity.Post{
		ID:              row.ID,
		AuthorID:        row.AuthorID,
		Text:            row.Text,
		Title:           row.Title,
		CommentsEnabled: row.CommentsEnabled,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}, nil
}

func (r *postRepository) GetPostsByAuthorID(ctx context.Context, authorID int64, limit, offset int) ([]*entity.Post, error) {
	rows, err := r.queries.GetPostsByAuthor(ctx, postssqlc.GetPostsByAuthorParams{
		AuthorID: authorID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrPostNotFound
		}
		return nil, err
	}

	listPosts := make([]*entity.Post, len(rows))
	for i, row := range rows {
		listPosts[i] = &entity.Post{
			ID:              row.ID,
			AuthorID:        row.AuthorID,
			Text:            row.Text,
			Title:           row.Title,
			CommentsEnabled: row.CommentsEnabled,
			CreatedAt:       row.CreatedAt.Time,
			UpdatedAt:       row.UpdatedAt.Time,
		}
	}
	return listPosts, nil
}

func (r *postRepository) GetPosts(ctx context.Context, limit, offset int) ([]*entity.Post, error) {
	rows, err := r.queries.GetPosts(ctx, postssqlc.GetPostsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	listPosts := make([]*entity.Post, len(rows))
	for i, row := range rows {
		listPosts[i] = &entity.Post{
			ID:              row.ID,
			AuthorID:        row.AuthorID,
			Text:            row.Text,
			Title:           row.Title,
			CommentsEnabled: row.CommentsEnabled,
			CreatedAt:       row.CreatedAt.Time,
			UpdatedAt:       row.UpdatedAt.Time,
		}
	}
	return listPosts, nil
}

func (r *postRepository) CreatePost(ctx context.Context, authorID int64, title string, text string, commentsEnabled bool) (*entity.Post, error) {
	row, err := r.queries.AddPost(ctx, postssqlc.AddPostParams{
		AuthorID:        authorID,
		Title:           title,
		Text:            text,
		CommentsEnabled: commentsEnabled,
	})
	if err != nil {
		return nil, err
	}
	return &entity.Post{
		ID:              row.ID,
		AuthorID:        row.AuthorID,
		Text:            row.Text,
		Title:           row.Title,
		CommentsEnabled: row.CommentsEnabled,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}, nil
}

func (r *postRepository) SetPostCommentsEnabled(ctx context.Context, id int64, enabled bool) (*entity.Post, error) {
	row, err := r.queries.SetCommentsEnabled(ctx, postssqlc.SetCommentsEnabledParams{
		ID:              id,
		CommentsEnabled: enabled,
	})
	if err != nil {
		return nil, err
	}
	return &entity.Post{
		ID:              row.ID,
		AuthorID:        row.AuthorID,
		Text:            row.Text,
		Title:           row.Title,
		CommentsEnabled: row.CommentsEnabled,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}, nil
}
