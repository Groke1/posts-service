package inmemory

import (
	"context"
	"testing"
	"time"

	"posts-service/internal/entity"

	"github.com/stretchr/testify/require"
)

func TestPostRepository_CreatePost(t *testing.T) {
	repo := NewPostRepository()

	post, err := repo.CreatePost(
		context.Background(),
		42,
		"Post title",
		"Post text",
		true,
	)

	require.NoError(t, err)
	require.Equal(t, int64(1), post.ID)
	require.Equal(t, int64(42), post.AuthorID)
	require.Equal(t, "Post title", post.Title)
	require.Equal(t, "Post text", post.Text)
	require.True(t, post.CommentsEnabled)
}

type getPostFunc func(*postRepository, context.Context, int64) (*entity.Post, error)

func runGetPostTests(t *testing.T, fn getPostFunc) {
	t.Helper()

	testCases := []struct {
		name         string
		setupData    setupDataFunc
		id           int64
		expectedPost *entity.Post
		expectedErr  error
	}{
		{
			name: "success",
			setupData: func(deps testDeps) {
				createPost(t, deps.repo, 42, "Post title", "Post text", true)
			},
			id: 1,
			expectedPost: &entity.Post{
				ID:              1,
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: true,
			},
		},
		{
			name:        "post not found",
			setupData:   func(deps testDeps) {},
			id:          1,
			expectedErr: entity.ErrPostNotFound,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps()
			tc.setupData(deps)

			post, err := fn(deps.repo, context.Background(), tc.id)

			checkError(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				return
			}

			requirePostEqual(t, tc.expectedPost, post)
		})
	}
}

func TestPostRepository_GetPostByID(t *testing.T) {
	t.Parallel()

	runGetPostTests(t,
		func(repo *postRepository, ctx context.Context, id int64) (*entity.Post, error) {
			return repo.GetPostByID(ctx, id)
		},
	)
}

func TestPostRepository_GetPostByIDForUpdate(t *testing.T) {
	t.Parallel()

	runGetPostTests(t,
		func(repo *postRepository, ctx context.Context, id int64) (*entity.Post, error) {
			return repo.GetPostByIDForUpdate(ctx, id)
		},
	)
}

func TestPostRepository_GetPosts(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		limit       int
		offset      int
		expectedIDs []int64
	}{
		{
			name:        "success",
			limit:       10,
			offset:      0,
			expectedIDs: []int64{3, 2, 1},
		},
		{
			name:        "success with limit",
			limit:       2,
			offset:      0,
			expectedIDs: []int64{3, 2},
		},
		{
			name:        "success with offset",
			limit:       10,
			offset:      1,
			expectedIDs: []int64{2, 1},
		},
		{
			name:        "offset out of range",
			limit:       10,
			offset:      100,
			expectedIDs: []int64{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := NewPostRepository()
			createPost(t, repo, 42, "Post 1", "Text 1", true)
			time.Sleep(time.Millisecond)
			createPost(t, repo, 43, "Post 2", "Text 2", true)
			time.Sleep(time.Millisecond)
			createPost(t, repo, 44, "Post 3", "Text 3", true)

			posts, err := repo.GetPosts(context.Background(), tc.limit, tc.offset)

			require.NoError(t, err)
			requirePostIDs(t, tc.expectedIDs, posts)
		})
	}
}

func TestPostRepository_GetPostsByAuthorID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		authorID    int64
		limit       int
		offset      int
		expectedIDs []int64
	}{
		{
			name:        "success",
			authorID:    42,
			limit:       10,
			offset:      0,
			expectedIDs: []int64{3, 1},
		},
		{
			name:        "success with limit",
			authorID:    42,
			limit:       1,
			offset:      0,
			expectedIDs: []int64{3},
		},
		{
			name:        "success with offset",
			authorID:    42,
			limit:       10,
			offset:      1,
			expectedIDs: []int64{1},
		},
		{
			name:        "author has no posts",
			authorID:    999,
			limit:       10,
			offset:      0,
			expectedIDs: []int64{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := NewPostRepository()

			createPost(t, repo, 42, "Post 1", "Text 1", true)
			time.Sleep(time.Millisecond)
			createPost(t, repo, 43, "Post 2", "Text 2", true)
			time.Sleep(time.Millisecond)
			createPost(t, repo, 42, "Post 3", "Text 3", true)

			posts, err := repo.GetPostsByAuthorID(context.Background(), tc.authorID,
				tc.limit, tc.offset)

			require.NoError(t, err)
			requirePostIDs(t, tc.expectedIDs, posts)
		})
	}
}

func TestPostRepository_SetPostCommentsEnabled(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		setupData    setupDataFunc
		id           int64
		enabled      bool
		expectedPost *entity.Post
		expectedErr  error
	}{
		{
			name: "success",
			setupData: func(deps testDeps) {
				createPost(t, deps.repo, 42, "Post title", "Post text", true)
			},
			id:      1,
			enabled: false,
			expectedPost: &entity.Post{
				ID:              1,
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: false,
			},
			expectedErr: nil,
		},
		{
			name:         "post not found",
			setupData:    func(deps testDeps) {},
			id:           1,
			enabled:      false,
			expectedPost: nil,
			expectedErr:  entity.ErrPostNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps()
			tc.setupData(deps)

			post, err := deps.repo.SetPostCommentsEnabled(
				context.Background(),
				tc.id,
				tc.enabled,
			)

			checkError(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				return
			}

			requirePostEqual(t, tc.expectedPost, post)
			require.False(t, post.UpdatedAt.IsZero())
		})
	}
}

type setupDataFunc func(deps testDeps)

type testDeps struct {
	repo *postRepository
}

func newTestDeps() testDeps {
	return testDeps{
		repo: NewPostRepository(),
	}
}

func createPost(t *testing.T, repo *postRepository, authorID int64, title string,
	text string, commentsEnabled bool) *entity.Post {
	t.Helper()

	post, err := repo.CreatePost(
		context.Background(),
		authorID,
		title,
		text,
		commentsEnabled,
	)
	require.NoError(t, err)
	require.NotNil(t, post)

	return post
}

func requirePostEqual(
	t *testing.T,
	expected *entity.Post,
	actual *entity.Post,
) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.AuthorID, actual.AuthorID)
	require.Equal(t, expected.Title, actual.Title)
	require.Equal(t, expected.Text, actual.Text)
	require.Equal(t, expected.CommentsEnabled, actual.CommentsEnabled)
}

func requirePostIDs(t *testing.T, expectedIDs []int64, posts []*entity.Post) {
	t.Helper()

	require.Equal(t, len(posts), len(expectedIDs))

	for i, post := range posts {
		require.Equal(t, expectedIDs[i], post.ID)
	}
}

func checkError(t *testing.T, actualErr error, expectedErr error) {
	t.Helper()

	if expectedErr == nil {
		require.NoError(t, actualErr)
		return
	}

	require.Error(t, actualErr)
	require.ErrorIs(t, actualErr, expectedErr)
}
