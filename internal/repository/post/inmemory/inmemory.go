package inmemory

import (
	"context"
	"posts-service/internal/entity"
	"sort"
	"sync"
	"time"
)

type postRepository struct {
	mx     sync.RWMutex
	postID int64
	posts  map[int64]*entity.Post
}

func NewPostRepository() *postRepository {
	return &postRepository{
		postID: 1,
		posts:  make(map[int64]*entity.Post),
	}
}

func (r *postRepository) GetPostByID(ctx context.Context, id int64) (*entity.Post, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	post, ok := r.posts[id]
	if !ok {
		return nil, entity.ErrPostNotFound
	}
	return post, nil
}

func (r *postRepository) GetPostByIDForUpdate(ctx context.Context, id int64) (*entity.Post, error) {
	return r.GetPostByID(ctx, id)
}

func (r *postRepository) GetPostsByAuthorID(ctx context.Context, authorID int64, limit, offset int) ([]*entity.Post, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	var postsByAuthor []*entity.Post
	for _, post := range r.posts {
		if post.AuthorID == authorID {
			postsByAuthor = append(postsByAuthor, post)
		}
	}
	return r.postsWithOffsetAndLimit(postsByAuthor, limit, offset), nil
}

func (r *postRepository) GetPosts(ctx context.Context, limit, offset int) ([]*entity.Post, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	var posts []*entity.Post
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return r.postsWithOffsetAndLimit(posts, limit, offset), nil
}

func (r *postRepository) CreatePost(ctx context.Context, authorID int64, title string, text string, commentsEnabled bool) (*entity.Post, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	t := time.Now()
	post := &entity.Post{
		ID:              r.postID,
		AuthorID:        authorID,
		Title:           title,
		Text:            text,
		CommentsEnabled: commentsEnabled,
		CreatedAt:       t,
		UpdatedAt:       t,
	}
	r.posts[post.ID] = post
	r.postID++
	return post, nil
}

func (r *postRepository) SetPostCommentsEnabled(ctx context.Context, id int64, enabled bool) (*entity.Post, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	post, ok := r.posts[id]
	if !ok {
		return nil, entity.ErrPostNotFound
	}
	post.CommentsEnabled = enabled
	post.UpdatedAt = time.Now()
	return post, nil
}

func (r *postRepository) postsWithOffsetAndLimit(posts []*entity.Post, limit, offset int) []*entity.Post {
	if offset >= len(posts) {
		return []*entity.Post{}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	end := min(offset+limit, len(posts))
	return posts[offset:end]
}
