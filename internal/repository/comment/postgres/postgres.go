package postgres

import (
	"context"
	"errors"
	"posts-service/internal/entity"
	commentssqlc "posts-service/internal/repository/comment/postgres/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type commentRepository struct {
	queries *commentssqlc.Queries
}

func NewCommentRepository(qdb commentssqlc.DBTX) *commentRepository {
	return &commentRepository{
		queries: commentssqlc.New(qdb),
	}
}

func (r *commentRepository) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	row, err := r.queries.GetComment(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrCommentNotFound
		}
		return nil, err
	}
	return &entity.Comment{
		ID:        row.ID,
		AuthorID:  row.AuthorID,
		PostID:    row.PostID,
		ParentID:  fromPgTypeInt8(row.ParentID),
		Text:      row.Text,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *commentRepository) CreateComment(ctx context.Context, postID, authorID int64, parentID *int64, text string) (*entity.Comment, error) {
	row, err := r.queries.AddComment(ctx, commentssqlc.AddCommentParams{
		PostID:   postID,
		AuthorID: authorID,
		Text:     text,
		ParentID: toPgTypeInt8(parentID),
	})
	if err != nil {
		return nil, err
	}

	return &entity.Comment{
		ID:        row.ID,
		AuthorID:  row.AuthorID,
		PostID:    row.PostID,
		ParentID:  fromPgTypeInt8(row.ParentID),
		Text:      row.Text,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *commentRepository) GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int) ([]*entity.Comment, error) {
	rows, err := r.queries.GetComments(ctx, commentssqlc.GetCommentsParams{
		PostID:   postID,
		ParentID: toPgTypeInt8(parentID),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}
	var comments []*entity.Comment
	for _, row := range rows {
		comments = append(comments, &entity.Comment{
			ID:        row.ID,
			AuthorID:  row.AuthorID,
			PostID:    row.PostID,
			ParentID:  fromPgTypeInt8(row.ParentID),
			Text:      row.Text,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}
	return comments, nil
}

func fromPgTypeInt8(id pgtype.Int8) *int64 {
	if !id.Valid {
		return nil
	}
	return &id.Int64
}

func toPgTypeInt8(id *int64) pgtype.Int8 {
	if id == nil {
		return pgtype.Int8{
			Valid: false,
		}
	}
	return pgtype.Int8{
		Valid: true,
		Int64: *id,
	}
}
