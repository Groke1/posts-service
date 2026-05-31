package entity

import (
	"errors"
	"time"
)

var (
	ErrCommentNotFound      = errors.New("comment not found")
	ErrCommentTooLong       = errors.New("comment too long")
	ErrCommentsDisabled     = errors.New("comments disabled")
	ErrEmptyCommentText     = errors.New("empty comment text")
	ErrParentCommentInvalid = errors.New("parent comment invalid")
)

type Comment struct {
	ID        int64
	ParentID  *int64
	PostID    int64
	AuthorID  int64
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
