package post

import (
	"context"
	"errors"
	"testing"

	"posts-service/internal/entity"
	"posts-service/internal/usecase/post/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var errRepo = errors.New("repo error")

func TestPostService_CreatePost(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		post         *entity.Post
		setupMocks   mocksSetupFunc
		expectedPost *entity.Post
		expectedErr  error
	}{
		{
			name: "success",
			post: &entity.Post{
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: true,
			},
			setupMocks: func(deps testDeps) {
				expectedPost := &entity.Post{
					ID:              100,
					AuthorID:        42,
					Title:           "Post title",
					Text:            "Post text",
					CommentsEnabled: true,
				}

				expectCreatePost(
					deps.postRepository,
					42,
					"Post title",
					"Post text",
					true,
					expectedPost,
					nil,
				)
			},
			expectedPost: &entity.Post{
				ID:              100,
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: true,
			},
			expectedErr: nil,
		},
		{
			name: "empty title",
			post: &entity.Post{
				AuthorID:        42,
				Title:           "",
				Text:            "Post text",
				CommentsEnabled: true,
			},
			setupMocks:  func(deps testDeps) {},
			expectedErr: entity.ErrPostTitleEmpty,
		},
		{
			name: "empty text",
			post: &entity.Post{
				AuthorID:        42,
				Title:           "Post title",
				Text:            "",
				CommentsEnabled: true,
			},
			setupMocks:  func(deps testDeps) {},
			expectedErr: entity.ErrPostTextEmpty,
		},
		{
			name: "repository error",
			post: &entity.Post{
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: true,
			},
			setupMocks: func(deps testDeps) {
				expectCreatePost(
					deps.postRepository,
					42,
					"Post title",
					"Post text",
					true,
					nil,
					errRepo,
				)
			},
			expectedPost: nil,
			expectedErr:  errRepo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			post, err := deps.service.CreatePost(context.Background(), tc.post)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedPost, post)
		})
	}
}

func TestPostService_GetPost(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		id           int64
		setupMocks   mocksSetupFunc
		expectedPost *entity.Post
		expectedErr  error
	}{
		{
			name: "success",
			id:   100,
			setupMocks: func(deps testDeps) {
				expectGetPostByID(deps.postRepository, 100, &entity.Post{
					ID:              100,
					AuthorID:        42,
					Title:           "Post title",
					Text:            "Post text",
					CommentsEnabled: true,
				}, nil)
			},
			expectedPost: &entity.Post{
				ID:              100,
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: true,
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			id:   100,
			setupMocks: func(deps testDeps) {
				expectGetPostByID(deps.postRepository, 100, nil, errRepo)
			},
			expectedErr: errRepo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			post, err := deps.service.GetPost(context.Background(), tc.id)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedPost, post)
		})
	}
}

func TestPostService_GetPosts(t *testing.T) {
	t.Parallel()

	authorID := int64(42)

	testCases := []struct {
		name          string
		authorID      *int64
		limit         int
		offset        int
		setupMocks    mocksSetupFunc
		expectedPosts []*entity.Post
		expectedErr   error
	}{
		{
			name:     "success with nil author id",
			authorID: nil,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPosts(deps.postRepository, 10, 0, []*entity.Post{
					{
						ID:              100,
						AuthorID:        42,
						Title:           "Post title 1",
						Text:            "Post text 1",
						CommentsEnabled: true,
					},
					{
						ID:              101,
						AuthorID:        43,
						Title:           "Post title 2",
						Text:            "Post text 2",
						CommentsEnabled: false,
					},
				}, nil)
			},
			expectedPosts: []*entity.Post{
				{
					ID:              100,
					AuthorID:        42,
					Title:           "Post title 1",
					Text:            "Post text 1",
					CommentsEnabled: true,
				},
				{
					ID:              101,
					AuthorID:        43,
					Title:           "Post title 2",
					Text:            "Post text 2",
					CommentsEnabled: false,
				},
			},
			expectedErr: nil,
		},
		{
			name:     "success with author id",
			authorID: &authorID,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPostsByAuthorID(deps.postRepository, 42, 10, 0, []*entity.Post{
					{
						ID:              100,
						AuthorID:        42,
						Title:           "Post title",
						Text:            "Post text",
						CommentsEnabled: true,
					},
				}, nil)
			},
			expectedPosts: []*entity.Post{
				{
					ID:              100,
					AuthorID:        42,
					Title:           "Post title",
					Text:            "Post text",
					CommentsEnabled: true,
				},
			},
			expectedErr: nil,
		},
		{
			name:        "invalid limit",
			authorID:    nil,
			limit:       0,
			offset:      0,
			setupMocks:  func(deps testDeps) {},
			expectedErr: entity.ErrInvalidPagination,
		},
		{
			name:        "invalid offset",
			authorID:    nil,
			limit:       10,
			offset:      -1,
			setupMocks:  func(deps testDeps) {},
			expectedErr: entity.ErrInvalidPagination,
		},
		{
			name:     "repository error with nil author id",
			authorID: nil,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPosts(deps.postRepository, 10, 0, nil, errRepo)
			},
			expectedErr: errRepo,
		},
		{
			name:     "repository error with author id",
			authorID: &authorID,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPostsByAuthorID(deps.postRepository, 42, 10, 0, nil, errRepo)
			},
			expectedErr: errRepo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			posts, err := deps.service.GetPosts(
				context.Background(),
				tc.authorID,
				tc.limit,
				tc.offset,
			)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedPosts, posts)
		})
	}
}

func TestPostService_SetPostCommentsEnabled(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		id           int64
		authorID     int64
		enabled      bool
		setupMocks   mocksSetupFunc
		expectedPost *entity.Post
		expectedErr  error
	}{
		{
			name:     "success",
			id:       100,
			authorID: 42,
			enabled:  false,
			setupMocks: func(deps testDeps) {
				expectGetPostByID(deps.postRepository, 100, &entity.Post{
					ID:              100,
					AuthorID:        42,
					Title:           "Post title",
					Text:            "Post text",
					CommentsEnabled: true,
				}, nil)

				expectSetPostCommentsEnabled(deps.postRepository, 100, false, &entity.Post{
					ID:              100,
					AuthorID:        42,
					Title:           "Post title",
					Text:            "Post text",
					CommentsEnabled: false,
				}, nil)
			},
			expectedPost: &entity.Post{
				ID:              100,
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: false,
			},
			expectedErr: nil,
		},
		{
			name:     "get post error",
			id:       100,
			authorID: 42,
			enabled:  false,
			setupMocks: func(deps testDeps) {
				expectGetPostByID(deps.postRepository, 100, nil, errRepo)
			},
			expectedPost: nil,
			expectedErr:  errRepo,
		},
		{
			name:     "forbidden",
			id:       100,
			authorID: 43,
			enabled:  false,
			setupMocks: func(deps testDeps) {
				expectGetPostByID(deps.postRepository, 100, &entity.Post{
					ID:              100,
					AuthorID:        42,
					Title:           "Post title",
					Text:            "Post text",
					CommentsEnabled: true,
				}, nil)
			},
			expectedErr: entity.ErrForbidden,
		},
		{
			name:     "set comments enabled error",
			id:       100,
			authorID: 42,
			enabled:  false,
			setupMocks: func(deps testDeps) {
				expectGetPostByID(deps.postRepository, 100, &entity.Post{
					ID:              100,
					AuthorID:        42,
					Title:           "Post title",
					Text:            "Post text",
					CommentsEnabled: true,
				}, nil)

				expectSetPostCommentsEnabled(deps.postRepository, 100, false, nil, errRepo)
			},
			expectedPost: nil,
			expectedErr:  errRepo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			post, err := deps.service.SetPostCommentsEnabled(
				context.Background(),
				tc.id,
				tc.authorID,
				tc.enabled,
			)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedPost, post)
		})
	}
}

type mocksSetupFunc func(deps testDeps)

type testDeps struct {
	postRepository *mocks.MockPostRepository
	service        *postService
}

func newTestDeps(t *testing.T) testDeps {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	postRepository := mocks.NewMockPostRepository(ctrl)

	return testDeps{
		postRepository: postRepository,
		service:        NewPostService(postRepository),
	}
}

func expectCreatePost(
	postRepository *mocks.MockPostRepository,
	authorID int64,
	title string,
	text string,
	commentsEnabled bool,
	post *entity.Post,
	err error,
) {
	postRepository.EXPECT().
		CreatePost(gomock.Any(), authorID, title, text, commentsEnabled).
		Return(post, err).Times(1)
}

func expectGetPostByID(
	postRepository *mocks.MockPostRepository,
	id int64,
	post *entity.Post,
	err error,
) {
	postRepository.EXPECT().
		GetPostByID(gomock.Any(), id).
		Return(post, err).Times(1)
}

func expectGetPostsByAuthorID(
	postRepository *mocks.MockPostRepository,
	authorID int64,
	limit int,
	offset int,
	posts []*entity.Post,
	err error,
) {
	postRepository.EXPECT().
		GetPostsByAuthorID(gomock.Any(), authorID, limit, offset).
		Return(posts, err).Times(1)
}

func expectGetPosts(
	postRepository *mocks.MockPostRepository,
	limit int,
	offset int,
	posts []*entity.Post,
	err error,
) {
	postRepository.EXPECT().
		GetPosts(gomock.Any(), limit, offset).
		Return(posts, err).Times(1)
}

func expectSetPostCommentsEnabled(
	postRepository *mocks.MockPostRepository,
	id int64,
	enabled bool,
	post *entity.Post,
	err error,
) {
	postRepository.EXPECT().
		SetPostCommentsEnabled(gomock.Any(), id, enabled).
		Return(post, err).Times(1)
}

func checkError(t *testing.T, actualErr error, expectedErr error) {
	t.Helper()

	if expectedErr == nil {
		require.NoError(t, actualErr)
		return
	}

	require.Error(t, actualErr)
	require.True(t, errors.Is(actualErr, expectedErr))
}
