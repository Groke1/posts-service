package comment

import (
	"context"
	"errors"
	"strings"
	"testing"

	"posts-service/internal/entity"
	"posts-service/internal/usecase/comment/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var errRepo = errors.New("repo error")
var errTx = errors.New("tx error")

func TestCommentService_CreateComment(t *testing.T) {
	t.Parallel()

	parentID := int64(10)

	testCases := []struct {
		name            string
		comment         *entity.Comment
		setupMocks      mocksSetupFunc
		expectedComment *entity.Comment
		expectedErr     error
	}{
		{
			name: "success with nil parent",
			comment: &entity.Comment{
				PostID:   100,
				AuthorID: 42,
				Text:     "Comment text",
			},
			setupMocks: func(deps testDeps) {
				expectedComment := &entity.Comment{
					ID:       1,
					PostID:   100,
					AuthorID: 42,
					Text:     "Comment text",
				}

				expectWithTx(deps.transactor, nil)
				expectGetPostByIDForUpdate(deps.postRepository, 100, &entity.Post{
					ID:              100,
					CommentsEnabled: true,
				}, nil)
				expectCreateComment(
					deps.commentRepository,
					100,
					42,
					nil,
					"Comment text",
					expectedComment,
					nil,
				)
				expectPublish(deps.publisher, expectedComment)
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
			name: "success with parent",
			comment: &entity.Comment{
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

				expectWithTx(deps.transactor, nil)
				expectGetPostByIDForUpdate(deps.postRepository, 100, &entity.Post{
					ID:              100,
					CommentsEnabled: true,
				}, nil)
				expectGetCommentByID(deps.commentRepository, 10, &entity.Comment{
					ID:     10,
					PostID: 100,
				}, nil)
				expectCreateComment(
					deps.commentRepository,
					100,
					42,
					&parentID,
					"Reply text",
					expectedComment,
					nil,
				)
				expectPublish(deps.publisher, expectedComment)
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
			name: "empty text",
			comment: &entity.Comment{
				PostID:   100,
				AuthorID: 42,
				Text:     "",
			},
			setupMocks:  func(deps testDeps) {},
			expectedErr: entity.ErrEmptyCommentText,
		},
		{
			name: "text too long",
			comment: &entity.Comment{
				PostID:   100,
				AuthorID: 42,
				Text:     strings.Repeat("a", 2001),
			},
			setupMocks:  func(deps testDeps) {},
			expectedErr: entity.ErrCommentTooLong,
		},
		{
			name: "transaction error",
			comment: &entity.Comment{
				PostID:   100,
				AuthorID: 42,
				Text:     "Comment text",
			},
			setupMocks: func(deps testDeps) {
				expectWithTx(deps.transactor, errTx)
			},
			expectedErr: errTx,
		},
		{
			name: "get post error",
			comment: &entity.Comment{
				PostID:   100,
				AuthorID: 42,
				Text:     "Comment text",
			},
			setupMocks: func(deps testDeps) {
				expectWithTx(deps.transactor, nil)
				expectGetPostByIDForUpdate(deps.postRepository, 100, nil, errRepo)
			},
			expectedComment: nil,
			expectedErr:     errRepo,
		},
		{
			name: "comments disabled",
			comment: &entity.Comment{
				PostID:   100,
				AuthorID: 42,
				Text:     "Comment text",
			},
			setupMocks: func(deps testDeps) {
				expectWithTx(deps.transactor, nil)
				expectGetPostByIDForUpdate(deps.postRepository, 100, &entity.Post{
					ID:              100,
					CommentsEnabled: false,
				}, nil)
			},
			expectedComment: nil,
			expectedErr:     entity.ErrCommentsDisabled,
		},
		{
			name: "get parent comment error",
			comment: &entity.Comment{
				PostID:   100,
				AuthorID: 42,
				ParentID: &parentID,
				Text:     "Reply text",
			},
			setupMocks: func(deps testDeps) {
				expectWithTx(deps.transactor, nil)
				expectGetPostByIDForUpdate(deps.postRepository, 100, &entity.Post{
					ID:              100,
					CommentsEnabled: true,
				}, nil)
				expectGetCommentByID(deps.commentRepository, 10, nil, errRepo)
			},
			expectedComment: nil,
			expectedErr:     errRepo,
		},
		{
			name: "parent comment invalid",
			comment: &entity.Comment{
				PostID:   100,
				AuthorID: 42,
				ParentID: &parentID,
				Text:     "Reply text",
			},
			setupMocks: func(deps testDeps) {
				expectWithTx(deps.transactor, nil)
				expectGetPostByIDForUpdate(deps.postRepository, 100, &entity.Post{
					ID:              100,
					CommentsEnabled: true,
				}, nil)
				expectGetCommentByID(deps.commentRepository, 10, &entity.Comment{
					ID:     10,
					PostID: 200,
				}, nil)
			},
			expectedComment: nil,
			expectedErr:     entity.ErrParentCommentInvalid,
		},
		{
			name: "create comment error",
			comment: &entity.Comment{
				PostID:   100,
				AuthorID: 42,
				Text:     "Comment text",
			},
			setupMocks: func(deps testDeps) {
				expectWithTx(deps.transactor, nil)
				expectGetPostByIDForUpdate(deps.postRepository, 100, &entity.Post{
					ID:              100,
					CommentsEnabled: true,
				}, nil)
				expectCreateComment(
					deps.commentRepository,
					100,
					42,
					nil,
					"Comment text",
					nil,
					errRepo,
				)
			},
			expectedErr: errRepo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			comment, err := deps.service.CreateComment(context.Background(), tc.comment)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedComment, comment)
		})
	}
}

func TestCommentService_GetComment(t *testing.T) {
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
				expectGetCommentByID(deps.commentRepository, 1, &entity.Comment{
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
			name: "repository error",
			id:   1,
			setupMocks: func(deps testDeps) {
				expectGetCommentByID(deps.commentRepository, 1, nil, errRepo)
			},
			expectedComment: nil,
			expectedErr:     errRepo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			comment, err := deps.service.GetComment(context.Background(), tc.id)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedComment, comment)
		})
	}
}

func TestCommentService_GetComments(t *testing.T) {
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
			name:     "success with nil parent",
			postID:   100,
			parentID: nil,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectedComments := []*entity.Comment{
					{
						ID:       1,
						PostID:   100,
						AuthorID: 42,
						Text:     "Comment text",
					},
				}

				expectGetPostByID(deps.postRepository, 100, &entity.Post{
					ID: 100,
				}, nil)
				expectGetComments(deps.commentRepository, 100, nil, 10, 0, expectedComments, nil)
			},
			expectedComments: []*entity.Comment{
				{
					ID:       1,
					PostID:   100,
					AuthorID: 42,
					Text:     "Comment text",
				},
			},
			expectedErr: nil,
		},
		{
			name:     "success with parent",
			postID:   100,
			parentID: &parentID,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectedComments := []*entity.Comment{
					{
						ID:       2,
						PostID:   100,
						AuthorID: 43,
						ParentID: &parentID,
						Text:     "Reply text",
					},
				}

				expectGetPostByID(deps.postRepository, 100, &entity.Post{
					ID: 100,
				}, nil)
				expectGetCommentByID(deps.commentRepository, 10, &entity.Comment{
					ID:     10,
					PostID: 100,
				}, nil)
				expectGetComments(deps.commentRepository, 100, &parentID, 10, 0, expectedComments, nil)
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
			name:        "invalid limit",
			postID:      100,
			parentID:    nil,
			limit:       0,
			offset:      0,
			setupMocks:  func(deps testDeps) {},
			expectedErr: entity.ErrInvalidPagination,
		},
		{
			name:        "invalid offset",
			postID:      100,
			parentID:    nil,
			limit:       10,
			offset:      -1,
			setupMocks:  func(deps testDeps) {},
			expectedErr: entity.ErrInvalidPagination,
		},
		{
			name:     "get post error",
			postID:   100,
			parentID: nil,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPostByID(deps.postRepository, 100, nil, errRepo)
			},
			expectedErr: errRepo,
		},
		{
			name:     "get parent comment error",
			postID:   100,
			parentID: &parentID,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPostByID(deps.postRepository, 100, &entity.Post{
					ID: 100,
				}, nil)
				expectGetCommentByID(deps.commentRepository, 10, nil, errRepo)
			},
			expectedErr: errRepo,
		},
		{
			name:     "parent comment invalid",
			postID:   100,
			parentID: &parentID,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPostByID(deps.postRepository, 100, &entity.Post{
					ID: 100,
				}, nil)
				expectGetCommentByID(deps.commentRepository, 10, &entity.Comment{
					ID:     10,
					PostID: 200,
				}, nil)
			},
			expectedErr: entity.ErrParentCommentInvalid,
		},
		{
			name:     "get comments error",
			postID:   100,
			parentID: nil,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPostByID(deps.postRepository, 100, &entity.Post{
					ID: 100,
				}, nil)
				expectGetComments(deps.commentRepository, 100, nil, 10, 0, nil, errRepo)
			},
			expectedErr: errRepo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			comments, err := deps.service.GetComments(
				context.Background(),
				tc.postID,
				tc.parentID,
				tc.limit,
				tc.offset,
			)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedComments, comments)
		})
	}
}

type mocksSetupFunc func(deps testDeps)

type testDeps struct {
	commentRepository *mocks.MockCommentRepository
	postRepository    *mocks.MockPostRepository
	transactor        *mocks.MockTransactor
	publisher         *mocks.Mockpublisher
	service           *commentService
}

func newTestDeps(t *testing.T) testDeps {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	commentRepository := mocks.NewMockCommentRepository(ctrl)
	postRepository := mocks.NewMockPostRepository(ctrl)
	transactor := mocks.NewMockTransactor(ctrl)
	publisher := mocks.NewMockpublisher(ctrl)

	return testDeps{
		commentRepository: commentRepository,
		postRepository:    postRepository,
		transactor:        transactor,
		publisher:         publisher,
		service:           NewCommentService(commentRepository, postRepository, transactor, publisher),
	}
}

func expectWithTx(
	transactor *mocks.MockTransactor,
	err error,
) {
	transactor.EXPECT().
		WithTx(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			if err != nil {
				return err
			}

			return fn(ctx)
		}).Times(1)
}

func expectCreateComment(
	commentRepository *mocks.MockCommentRepository,
	postID int64,
	authorID int64,
	parentID *int64,
	text string,
	comment *entity.Comment,
	err error,
) {
	commentRepository.EXPECT().
		CreateComment(gomock.Any(), postID, authorID, parentID, text).
		Return(comment, err).Times(1)
}

func expectGetCommentByID(
	commentRepository *mocks.MockCommentRepository,
	id int64,
	comment *entity.Comment,
	err error,
) {
	commentRepository.EXPECT().
		GetCommentByID(gomock.Any(), id).
		Return(comment, err).Times(1)
}

func expectGetComments(
	commentRepository *mocks.MockCommentRepository,
	postID int64,
	parentID *int64,
	limit int,
	offset int,
	comments []*entity.Comment,
	err error,
) {
	commentRepository.EXPECT().
		GetComments(gomock.Any(), postID, parentID, limit, offset).
		Return(comments, err).Times(1)
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

func expectGetPostByIDForUpdate(
	postRepository *mocks.MockPostRepository,
	id int64,
	post *entity.Post,
	err error,
) {
	postRepository.EXPECT().
		GetPostByIDForUpdate(gomock.Any(), id).
		Return(post, err).Times(1)
}

func expectPublish(
	publisher *mocks.Mockpublisher,
	comment *entity.Comment,
) {
	publisher.EXPECT().
		Publish(comment).
		Times(1)
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
