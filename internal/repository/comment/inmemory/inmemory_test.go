package inmemory

import (
	"context"
	"testing"
	"time"

	"posts-service/internal/entity"

	"github.com/stretchr/testify/require"
)

func TestCommentRepository_CreateComment(t *testing.T) {
	t.Parallel()

	parentID := int64(10)

	testCases := []struct {
		name     string
		postID   int64
		authorID int64
		parentID *int64
		text     string
	}{
		{
			name:     "success with nil parent",
			postID:   100,
			authorID: 42,
			parentID: nil,
			text:     "Comment text",
		},
		{
			name:     "success with parent",
			postID:   100,
			authorID: 42,
			parentID: &parentID,
			text:     "Reply text",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps()

			comment, err := deps.repo.CreateComment(
				context.Background(), tc.postID,
				tc.authorID, tc.parentID, tc.text,
			)

			require.NoError(t, err)
			require.NotNil(t, comment)
			require.Equal(t, int64(1), comment.ID)
			require.Equal(t, tc.postID, comment.PostID)
			require.Equal(t, tc.authorID, comment.AuthorID)
			require.Equal(t, tc.parentID, comment.ParentID)
			require.Equal(t, tc.text, comment.Text)
		})
	}
}

func TestCommentRepository_GetCommentByID(t *testing.T) {
	t.Parallel()

	parentID := int64(10)

	testCases := []struct {
		name            string
		setupData       setupDataFunc
		id              int64
		expectedComment *entity.Comment
		expectedErr     error
	}{
		{
			name: "success with nil parent",
			setupData: func(deps testDeps) {
				createComment(t, deps.repo, 100, 42, nil, "Comment text")
			},
			id: 1,
			expectedComment: &entity.Comment{
				ID:       1,
				PostID:   100,
				AuthorID: 42,
				ParentID: nil,
				Text:     "Comment text",
			},
			expectedErr: nil,
		},
		{
			name: "success with parent",
			setupData: func(deps testDeps) {
				createComment(t, deps.repo, 100, 42, &parentID, "Reply text")
			},
			id: 1,
			expectedComment: &entity.Comment{
				ID:       1,
				PostID:   100,
				AuthorID: 42,
				ParentID: &parentID,
				Text:     "Reply text",
			},
			expectedErr: nil,
		},
		{
			name:            "comment not found",
			setupData:       func(deps testDeps) {},
			id:              1,
			expectedComment: nil,
			expectedErr:     entity.ErrCommentNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps()
			tc.setupData(deps)

			comment, err := deps.repo.GetCommentByID(context.Background(), tc.id)

			checkError(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				return
			}

			requireCommentEqual(t, tc.expectedComment, comment)
		})
	}
}

func TestCommentRepository_GetComments(t *testing.T) {
	t.Parallel()

	parentID := int64(10)
	anotherParentID := int64(20)

	testCases := []struct {
		name        string
		postID      int64
		parentID    *int64
		limit       int
		offset      int
		setupData   setupDataFunc
		expectedIDs []int64
	}{
		{
			name:     "success",
			postID:   100,
			parentID: nil,
			limit:    10,
			offset:   0,
			setupData: func(deps testDeps) {
				createComment(t, deps.repo, 100, 42, nil, "Comment 1")
				time.Sleep(time.Millisecond)
				createComment(t, deps.repo, 100, 43, nil, "Comment 2")
				time.Sleep(time.Millisecond)
				createComment(t, deps.repo, 100, 44, nil, "Comment 3")
			},
			expectedIDs: []int64{1, 2, 3},
		},
		{
			name:     "filters by post id",
			postID:   100,
			parentID: nil,
			limit:    10,
			offset:   0,
			setupData: func(deps testDeps) {
				createComment(t, deps.repo, 100, 42, nil, "Comment 1")
				time.Sleep(time.Nanosecond)
				createComment(t, deps.repo, 200, 43, nil, "Comment 2")
				time.Sleep(time.Nanosecond)
				createComment(t, deps.repo, 100, 44, nil, "Comment 3")
			},
			expectedIDs: []int64{1, 3},
		},
		{
			name:     "filters by parent id",
			postID:   100,
			parentID: &parentID,
			limit:    10,
			offset:   0,
			setupData: func(deps testDeps) {
				createComment(t, deps.repo, 100, 42, nil, "Root comment")
				time.Sleep(time.Millisecond)
				createComment(t, deps.repo, 100, 43, &parentID, "Reply 1")
				time.Sleep(time.Millisecond)
				createComment(t, deps.repo, 100, 44, &anotherParentID, "Reply 2")
				time.Sleep(time.Millisecond)
				createComment(t, deps.repo, 100, 45, &parentID, "Reply 3")
			},
			expectedIDs: []int64{2, 4},
		},
		{
			name:     "success with limit",
			postID:   100,
			parentID: nil,
			limit:    2,
			offset:   0,
			setupData: func(deps testDeps) {
				createComment(t, deps.repo, 100, 42, nil, "Comment 1")
				time.Sleep(time.Millisecond)
				createComment(t, deps.repo, 100, 43, nil, "Comment 2")
				time.Sleep(time.Millisecond)
				createComment(t, deps.repo, 100, 44, nil, "Comment 3")
			},
			expectedIDs: []int64{1, 2},
		},
		{
			name:     "success with offset",
			postID:   100,
			parentID: nil,
			limit:    10,
			offset:   1,
			setupData: func(deps testDeps) {
				createComment(t, deps.repo, 100, 42, nil, "Comment 1")
				time.Sleep(time.Millisecond)
				createComment(t, deps.repo, 100, 43, nil, "Comment 2")
				time.Sleep(time.Millisecond)
				createComment(t, deps.repo, 100, 44, nil, "Comment 3")
			},
			expectedIDs: []int64{2, 3},
		},
		{
			name:     "offset out of range",
			postID:   100,
			parentID: nil,
			limit:    10,
			offset:   1,
			setupData: func(deps testDeps) {
				createComment(t, deps.repo, 100, 42, nil, "Comment 1")
			},
			expectedIDs: []int64{},
		},
		{
			name:     "no comments for parent",
			postID:   100,
			parentID: &parentID,
			limit:    10,
			offset:   0,
			setupData: func(deps testDeps) {
				createComment(t, deps.repo, 100, 42, nil, "Root comment")
				createComment(t, deps.repo, 100, 43, &anotherParentID, "Reply")
			},
			expectedIDs: []int64{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps()
			tc.setupData(deps)

			comments, err := deps.repo.GetComments(
				context.Background(), tc.postID,
				tc.parentID, tc.limit,
				tc.offset,
			)

			require.NoError(t, err)
			requireCommentIDs(t, tc.expectedIDs, comments)
		})
	}
}

type setupDataFunc func(deps testDeps)

type testDeps struct {
	repo *commentRepository
}

func newTestDeps() testDeps {
	return testDeps{
		repo: NewCommentRepository(),
	}
}

func createComment(t *testing.T, repo *commentRepository, postID int64,
	authorID int64, parentID *int64, text string,
) *entity.Comment {
	t.Helper()

	comment, err := repo.CreateComment(
		context.Background(),
		postID,
		authorID,
		parentID,
		text,
	)
	require.NoError(t, err)
	require.NotNil(t, comment)

	return comment
}

func requireCommentEqual(
	t *testing.T,
	expected *entity.Comment,
	actual *entity.Comment,
) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.PostID, actual.PostID)
	require.Equal(t, expected.AuthorID, actual.AuthorID)
	require.Equal(t, expected.ParentID, actual.ParentID)
	require.Equal(t, expected.Text, actual.Text)
}

func requireCommentIDs(t *testing.T, expectedIDs []int64, comments []*entity.Comment) {
	t.Helper()

	require.Equal(t, len(comments), len(expectedIDs))
	for i, comment := range comments {
		require.Equal(t, expectedIDs[i], comment.ID)
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
