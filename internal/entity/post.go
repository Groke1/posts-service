package entity

import (
	"errors"
	"time"
)

var (
	ErrPostNotFound      = errors.New("post not found")
	ErrPostTextEmpty     = errors.New("post text is empty")
	ErrPostTitleEmpty    = errors.New("post title is empty")
	ErrInvalidPagination = errors.New("invalid pagination")
	ErrForbidden         = errors.New("forbidden")
)

type Post struct {
	ID              int64
	Title           string
	Text            string
	AuthorID        int64
	CommentsEnabled bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
