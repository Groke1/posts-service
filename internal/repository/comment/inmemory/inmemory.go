package inmemory

import (
	"context"
	"posts-service/internal/entity"
	"sort"
	"sync"
	"time"
)

type commentRepository struct {
	mx        sync.RWMutex
	commentID int64
	comments  map[int64]*entity.Comment
}

func NewCommentRepository() *commentRepository {
	return &commentRepository{
		commentID: 1,
		comments:  make(map[int64]*entity.Comment),
	}
}

func (r *commentRepository) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	comment, ok := r.comments[id]
	if !ok {
		return nil, entity.ErrCommentNotFound
	}
	return comment, nil
}

func (r *commentRepository) CreateComment(ctx context.Context, postID, authorID int64, parentID *int64, text string) (*entity.Comment, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	t := time.Now()
	comment := &entity.Comment{
		ID:        r.commentID,
		PostID:    postID,
		AuthorID:  authorID,
		ParentID:  parentID,
		Text:      text,
		CreatedAt: t,
		UpdatedAt: t,
	}
	r.comments[comment.ID] = comment
	r.commentID++
	return comment, nil
}

func (r *commentRepository) GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int) ([]*entity.Comment, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	comments := make([]*entity.Comment, 0)

	for _, comment := range r.comments {
		if comment.PostID != postID || !sameParentID(comment.ParentID, parentID) {
			continue
		}

		comments = append(comments, comment)
	}

	if offset >= len(comments) {
		return []*entity.Comment{}, nil
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})

	end := min(offset + limit, len(comments))

	return comments[offset:end], nil
}

func sameParentID(id1, id2 *int64) bool {
	if id1 == nil && id2 == nil {
		return true
	}

	if id1 == nil || id2 == nil {
		return false
	}

	return *id1 == *id2
}
