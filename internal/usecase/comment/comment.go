package comment

import (
	"context"
	"fmt"
	"posts-service/internal/entity"
	"unicode/utf8"
)

//go:generate mockgen -source=comment.go -destination=mocks/comment_mocks.go -package=mocks

type (
	CommentRepository interface {
		CreateComment(ctx context.Context, postID, authorID int64, parentID *int64, text string) (*entity.Comment, error)
		GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
		GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int) ([]*entity.Comment, error)
	}
	PostRepository interface {
		GetPostByID(ctx context.Context, id int64) (*entity.Post, error)
		GetPostByIDForUpdate(ctx context.Context, id int64) (*entity.Post, error)
	}

	Transactor interface {
		WithTx(ctx context.Context, fn func(ctx context.Context) error) error
	}

	publisher interface {
		Publish(comment *entity.Comment)
	}
)

type commentService struct {
	commentRepository CommentRepository
	postRepository    PostRepository
	transactor        Transactor
	publisher         publisher
}

func NewCommentService(
	commentRepository CommentRepository,
	postRepository PostRepository,
	transactor Transactor,
	publisher publisher) *commentService {
	return &commentService{
		commentRepository: commentRepository,
		postRepository:    postRepository,
		transactor:        transactor,
		publisher:         publisher,
	}
}

func (s *commentService) CreateComment(ctx context.Context, comment *entity.Comment) (resComment *entity.Comment, err error) {
	if utf8.RuneCountInString(comment.Text) > 2000 {
		return nil, entity.ErrCommentTooLong
	}
	if comment.Text == "" {
		return nil, entity.ErrEmptyCommentText
	}
	err = s.transactor.WithTx(ctx, func(ctx context.Context) error {
		post, err := s.postRepository.GetPostByIDForUpdate(ctx, comment.PostID)
		if err != nil {
			return fmt.Errorf("failed to get post by id %d: %w", comment.PostID, err)
		}
		if !post.CommentsEnabled {
			return entity.ErrCommentsDisabled
		}

		if comment.ParentID != nil {
			parentComment, err := s.commentRepository.GetCommentByID(ctx, *comment.ParentID)
			if err != nil {
				return err
			}

			if parentComment.PostID != comment.PostID {
				return entity.ErrParentCommentInvalid
			}
		}

		resComment, err = s.commentRepository.CreateComment(ctx, comment.PostID,
			comment.AuthorID, comment.ParentID, comment.Text)
		if err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.publisher.Publish(resComment)
	return resComment, nil
}

func (s *commentService) GetComment(ctx context.Context, id int64) (*entity.Comment, error) {
	return s.commentRepository.GetCommentByID(ctx, id)
}

func (s *commentService) GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int) ([]*entity.Comment, error) {
	if limit <= 0 || offset < 0 {
		return nil, entity.ErrInvalidPagination
	}

	_, err := s.postRepository.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if parentID != nil {
		parentComment, err := s.commentRepository.GetCommentByID(ctx, *parentID)
		if err != nil {
			return nil, err
		}

		if parentComment.PostID != postID {
			return nil, entity.ErrParentCommentInvalid
		}
	}

	return s.commentRepository.GetComments(ctx, postID, parentID, limit, offset)
}
