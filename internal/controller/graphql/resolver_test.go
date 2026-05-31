package graphql

import (
	"context"
	"errors"
	"testing"

	"posts-service/internal/controller/graphql/converter"
	"posts-service/internal/controller/graphql/mocks"
	"posts-service/internal/entity"
	"posts-service/pkg/generated/posts/graphql/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var errService = errors.New("service error")

func TestMutationResolver_CreatePost(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		input        model.CreatePostInput
		setupMocks   mocksSetupFunc
		expectedPost *entity.Post
		expectedErr  error
	}{
		{
			name: "success",
			input: model.CreatePostInput{
				AuthorID: 42,
				Title:    "Post title",
				Text:     "Post text",
			},
			setupMocks: func(deps testDeps) {
				expectedPost := &entity.Post{
					ID:              100,
					AuthorID:        42,
					Title:           "Post title",
					Text:            "Post text",
					CommentsEnabled: true,
				}

				expectCreatePost(t, deps.postService, &entity.Post{
					AuthorID:        42,
					Title:           "Post title",
					Text:            "Post text",
					CommentsEnabled: true,
				}, expectedPost, nil)
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
			name: "service error",
			input: model.CreatePostInput{
				AuthorID: 42,
				Title:    "Post title",
				Text:     "Post text",
			},
			setupMocks: func(deps testDeps) {
				expectCreatePost(t, deps.postService, &entity.Post{
					AuthorID:        42,
					Title:           "Post title",
					Text:            "Post text",
					CommentsEnabled: true,
				}, nil, errService)
			},
			expectedErr: errService,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			post, err := deps.mutationResolver.CreatePost(context.Background(), tc.input)

			checkError(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				require.Nil(t, post)
				return
			}

			require.Equal(t, converter.ToGraphQLPost(tc.expectedPost), post)
		})
	}
}

func TestMutationResolver_SetPostCommentsEnabled(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		postID       int64
		authorID     int64
		enabled      bool
		setupMocks   mocksSetupFunc
		expectedPost *entity.Post
		expectedErr  error
	}{
		{
			name:     "success",
			postID:   100,
			authorID: 42,
			enabled:  false,
			setupMocks: func(deps testDeps) {
				expectSetPostCommentsEnabled(deps.postService, 100, 42, false, &entity.Post{
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
			name:     "service error",
			postID:   100,
			authorID: 42,
			enabled:  false,
			setupMocks: func(deps testDeps) {
				expectSetPostCommentsEnabled(deps.postService, 100, 42, false, nil, errService)
			},
			expectedErr: errService,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			post, err := deps.mutationResolver.SetPostCommentsEnabled(
				context.Background(),
				tc.postID,
				tc.authorID,
				tc.enabled,
			)

			checkError(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				return
			}

			require.Equal(t, converter.ToGraphQLPost(tc.expectedPost), post)
		})
	}
}

func TestQueryResolver_Post(t *testing.T) {
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
				expectGetPost(deps.postService, 100, &entity.Post{
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
			name: "service error",
			id:   100,
			setupMocks: func(deps testDeps) {
				expectGetPost(deps.postService, 100, nil, errService)
			},
			expectedErr: errService,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			post, err := deps.queryResolver.Post(context.Background(), tc.id)

			checkError(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				return
			}

			require.Equal(t, converter.ToGraphQLPost(tc.expectedPost), post)
		})
	}
}

func TestQueryResolver_Posts(t *testing.T) {
	t.Parallel()

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
			name:     "success",
			authorID: nil,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPosts(deps.postService, nil, 10, 0, []*entity.Post{
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
			name:     "service error",
			authorID: nil,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPosts(deps.postService, nil, 10, 0, nil, errService)
			},
			expectedErr: errService,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			posts, err := deps.queryResolver.Posts(
				context.Background(),
				tc.authorID,
				tc.limit,
				tc.offset,
			)

			checkError(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				require.Nil(t, posts)
				return
			}

			require.Equal(t, converter.ToGraphQLPosts(tc.expectedPosts), posts)
		})
	}
}

func TestMutationResolver_CreateComment(t *testing.T) {
	t.Parallel()

	parentID := int64(10)

	testCases := []struct {
		name            string
		input           model.CreateCommentInput
		setupMocks      mocksSetupFunc
		expectedComment *entity.Comment
		expectedErr     error
	}{
		{
			name: "success",
			input: model.CreateCommentInput{
				PostID:   100,
				AuthorID: 42,
				ParentID: &parentID,
				Text:     "Reply text",
			},
			setupMocks: func(deps testDeps) {
				expectedComment := &entity.Comment{
					ID:       2,
					PostID:   100,
					AuthorID: 42,
					ParentID: &parentID,
					Text:     "Reply text",
				}

				expectCreateComment(t, deps.commentService, &entity.Comment{
					PostID:   100,
					AuthorID: 42,
					ParentID: &parentID,
					Text:     "Reply text",
				}, expectedComment, nil)
			},
			expectedComment: &entity.Comment{
				ID:       2,
				PostID:   100,
				AuthorID: 42,
				ParentID: &parentID,
				Text:     "Reply text",
			},
			expectedErr: nil,
		},
		{
			name: "service error",
			input: model.CreateCommentInput{
				PostID:   100,
				AuthorID: 42,
				Text:     "Comment text",
			},
			setupMocks: func(deps testDeps) {
				expectCreateComment(t, deps.commentService, &entity.Comment{
					PostID:   100,
					AuthorID: 42,
					Text:     "Comment text",
				}, nil, errService)
			},
			expectedErr: errService,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			comment, err := deps.mutationResolver.CreateComment(context.Background(), tc.input)

			checkError(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				return
			}

			require.Equal(t, converter.ToGraphQLComment(tc.expectedComment), comment)
		})
	}
}

func TestQueryResolver_Comment(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		id              int64
		setupMocks      mocksSetupFunc
		expectedComment *entity.Comment
		expectedErr     error
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(deps testDeps) {
				expectGetComment(deps.commentService, 1, &entity.Comment{
					ID:       1,
					PostID:   100,
					AuthorID: 42,
					Text:     "Comment text",
				}, nil)
			},
			expectedComment: &entity.Comment{
				ID:       1,
				PostID:   100,
				AuthorID: 42,
				Text:     "Comment text",
			},
			expectedErr: nil,
		},
		{
			name: "service error",
			id:   1,
			setupMocks: func(deps testDeps) {
				expectGetComment(deps.commentService, 1, nil, errService)
			},
			expectedComment: nil,
			expectedErr:     errService,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			comment, err := deps.queryResolver.Comment(context.Background(), tc.id)

			checkError(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				require.Nil(t, comment)
				return
			}

			require.Equal(t, converter.ToGraphQLComment(tc.expectedComment), comment)
		})
	}
}

func TestQueryResolver_Comments(t *testing.T) {
	t.Parallel()

	parentID := int64(10)

	testCases := []struct {
		name             string
		postID           int64
		parentID         *int64
		limit            int
		offset           int
		setupMocks       mocksSetupFunc
		expectedComments []*entity.Comment
		expectedErr      error
	}{
		{
			name:     "success",
			postID:   100,
			parentID: &parentID,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetComments(deps.commentService, 100, &parentID, 10, 0, []*entity.Comment{
					{
						ID:       2,
						PostID:   100,
						AuthorID: 43,
						ParentID: &parentID,
						Text:     "Reply text",
					},
				}, nil)
			},
			expectedComments: []*entity.Comment{
				{
					ID:       2,
					PostID:   100,
					AuthorID: 43,
					ParentID: &parentID,
					Text:     "Reply text",
				},
			},
			expectedErr: nil,
		},
		{
			name:     "service error",
			postID:   100,
			parentID: nil,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetComments(deps.commentService, 100, nil, 10, 0, nil, errService)
			},
			expectedErr: errService,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			comments, err := deps.queryResolver.Comments(
				context.Background(),
				tc.postID,
				tc.parentID,
				tc.limit,
				tc.offset,
			)

			checkError(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				return
			}

			require.Equal(t, converter.ToGraphQLComments(tc.expectedComments), comments)
		})
	}
}

type mocksSetupFunc func(deps testDeps)

type testDeps struct {
	ctx                  context.Context
	postService          *mocks.MockpostService
	commentService       *mocks.MockcommentService
	resolver             *Resolver
	mutationResolver     *mutationResolver
	queryResolver        *queryResolver
	subscriptionResolver *subscriptionResolver
}

func newTestDeps(t *testing.T) testDeps {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	postService := mocks.NewMockpostService(ctrl)
	commentService := mocks.NewMockcommentService(ctrl)
	commentSub := mocks.NewMockcommentSubscriber(ctrl)

	deps := testDeps{
		ctx:            context.Background(),
		postService:    postService,
		commentService: commentService,
	}

	deps.resolver = NewResolver(postService, commentService, commentSub)
	deps.mutationResolver = &mutationResolver{deps.resolver}
	deps.queryResolver = &queryResolver{deps.resolver}
	deps.subscriptionResolver = &subscriptionResolver{deps.resolver}

	return deps
}

func expectCreatePost(
	t *testing.T,
	postService *mocks.MockpostService,
	expectedPost *entity.Post,
	post *entity.Post,
	err error,
) {
	t.Helper()

	postService.EXPECT().
		CreatePost(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, actualPost *entity.Post) (*entity.Post, error) {
			require.Equal(t, expectedPost, actualPost)
			return post, err
		}).Times(1)
}

func expectSetPostCommentsEnabled(
	postService *mocks.MockpostService,
	id int64,
	authorID int64,
	enabled bool,
	post *entity.Post,
	err error,
) {
	postService.EXPECT().
		SetPostCommentsEnabled(gomock.Any(), id, authorID, enabled).
		Return(post, err).Times(1)
}

func expectGetPost(
	postService *mocks.MockpostService,
	id int64,
	post *entity.Post,
	err error,
) {
	postService.EXPECT().
		GetPost(gomock.Any(), id).
		Return(post, err).Times(1)
}

func expectGetPosts(
	postService *mocks.MockpostService,
	authorID *int64,
	limit int,
	offset int,
	posts []*entity.Post,
	err error,
) {
	postService.EXPECT().
		GetPosts(gomock.Any(), authorID, limit, offset).
		Return(posts, err).Times(1)
}

func expectCreateComment(
	t *testing.T,
	commentService *mocks.MockcommentService,
	expectedComment *entity.Comment,
	comment *entity.Comment,
	err error,
) {
	t.Helper()

	commentService.EXPECT().
		CreateComment(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, actualComment *entity.Comment) (*entity.Comment, error) {
			require.Equal(t, expectedComment, actualComment)
			return comment, err
		}).Times(1)
}

func expectGetComment(
	commentService *mocks.MockcommentService,
	id int64,
	comment *entity.Comment,
	err error,
) {
	commentService.EXPECT().
		GetComment(gomock.Any(), id).
		Return(comment, err).Times(1)
}

func expectGetComments(
	commentService *mocks.MockcommentService, postID int64,
	parentID *int64, limit int,
	offset int, comments []*entity.Comment,
	err error,
) {
	commentService.EXPECT().
		GetComments(gomock.Any(), postID, parentID, limit, offset).
		Return(comments, err).Times(1)
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
