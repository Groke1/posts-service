package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"posts-service/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

var errDB = errors.New("db error")

func TestCommentRepository_GetCommentByID(t *testing.T) {
	t.Parallel()

	now := time.Now()
	parentID := int64(10)

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
				expectGetComment(deps.dbMock, 1, commentRow{
					id:        1,
					postID:    100,
					authorID:  42,
					parentID:  &parentID,
					text:      "Reply text",
					createdAt: now,
					updatedAt: now,
				}, nil)
			},
			expectedComment: &entity.Comment{
				ID:        1,
				PostID:    100,
				AuthorID:  42,
				ParentID:  &parentID,
				Text:      "Reply text",
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectedErr: nil,
		},
		{
			name: "comment not found",
			id:   1,
			setupMocks: func(deps testDeps) {
				expectGetComment(deps.dbMock, 1, commentRow{}, pgx.ErrNoRows)
			},
			expectedComment: nil,
			expectedErr:     entity.ErrCommentNotFound,
		},
		{
			name: "db error",
			id:   1,
			setupMocks: func(deps testDeps) {
				expectGetComment(deps.dbMock, 1, commentRow{}, errDB)
			},
			expectedComment: nil,
			expectedErr:     errDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			comment, err := deps.repo.GetCommentByID(context.Background(), tc.id)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedComment, comment)
		})
	}
}

func TestCommentRepository_CreateComment(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testCases := []struct {
		name            string
		postID          int64
		authorID        int64
		parentID        *int64
		text            string
		setupMocks      mocksSetupFunc
		expectedComment *entity.Comment
		expectedErr     error
	}{
		{
			name:     "success",
			postID:   100,
			authorID: 42,
			parentID: nil,
			text:     "Comment text",
			setupMocks: func(deps testDeps) {
				expectAddComment(deps.dbMock, 100, 42, nil, "Comment text", commentRow{
					id:        1,
					postID:    100,
					authorID:  42,
					parentID:  nil,
					text:      "Comment text",
					createdAt: now,
					updatedAt: now,
				}, nil)
			},
			expectedComment: &entity.Comment{
				ID:        1,
				PostID:    100,
				AuthorID:  42,
				ParentID:  nil,
				Text:      "Comment text",
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectedErr: nil,
		},
		{
			name:     "db error",
			postID:   100,
			authorID: 42,
			parentID: nil,
			text:     "Comment text",
			setupMocks: func(deps testDeps) {
				expectAddComment(deps.dbMock, 100, 42, nil, "Comment text", commentRow{}, errDB)
			},
			expectedComment: nil,
			expectedErr:     errDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			comment, err := deps.repo.CreateComment(
				context.Background(), tc.postID,
				tc.authorID, tc.parentID,
				tc.text,
			)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedComment, comment)
		})
	}
}

func TestCommentRepository_GetComments(t *testing.T) {
	t.Parallel()

	now := time.Now()
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
				expectGetComments(deps.dbMock, 100, &parentID, 10, 0, []commentRow{
					{
						id:        3,
						postID:    100,
						authorID:  44,
						parentID:  &parentID,
						text:      "Reply text",
						createdAt: now,
						updatedAt: now,
					},
				}, nil)
			},
			expectedComments: []*entity.Comment{
				{
					ID:        3,
					PostID:    100,
					AuthorID:  44,
					ParentID:  &parentID,
					Text:      "Reply text",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			expectedErr: nil,
		},
		{
			name:     "db error",
			postID:   100,
			parentID: nil,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetComments(deps.dbMock, 100, nil, 10, 0, nil, errDB)
			},
			expectedComments: nil,
			expectedErr:      errDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			comments, err := deps.repo.GetComments(
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
	dbMock pgxmock.PgxConnIface
	repo   *commentRepository
}

type commentRow struct {
	id        int64
	postID    int64
	authorID  int64
	parentID  *int64
	text      string
	createdAt time.Time
	updatedAt time.Time
}

func newTestDeps(t *testing.T) testDeps {
	t.Helper()

	dbMock, err := pgxmock.NewConn()
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, dbMock.ExpectationsWereMet())
		dbMock.Close(context.Background())
	})

	return testDeps{
		dbMock: dbMock,
		repo:   NewCommentRepository(dbMock),
	}
}

func expectGetComment(
	dbMock pgxmock.PgxConnIface,
	id int64,
	row commentRow,
	err error,
) {
	query := dbMock.ExpectQuery(".*").
		WithArgs(id)

	if err != nil {
		query.WillReturnError(err)
		return
	}

	query.WillReturnRows(newCommentRows([]commentRow{row}))
}

func expectAddComment(
	dbMock pgxmock.PgxConnIface,
	postID int64,
	authorID int64,
	parentID *int64,
	text string,
	row commentRow,
	err error,
) {
	query := dbMock.ExpectQuery(".*").
		WithArgs(postID, authorID, toPgTypeInt8(parentID), text)

	if err != nil {
		query.WillReturnError(err)
		return
	}

	query.WillReturnRows(newCommentRows([]commentRow{row}))
}

func expectGetComments(
	dbMock pgxmock.PgxConnIface,
	postID int64,
	parentID *int64,
	limit int,
	offset int,
	rows []commentRow,
	err error,
) {
	query := dbMock.ExpectQuery(".*").
		WithArgs(postID, toPgTypeInt8(parentID), int32(offset), int32(limit))

	if err != nil {
		query.WillReturnError(err)
		return
	}

	query.WillReturnRows(newCommentRows(rows))
}

func newCommentRows(rows []commentRow) *pgxmock.Rows {
	result := pgxmock.NewRows([]string{
		"id", "post_id",
		"author_id", "parent_id",
		"text", "created_at",
		"updated_at",
	})

	for _, row := range rows {
		result.AddRow(
			row.id,
			row.postID,
			toPgTypeInt8(row.parentID),
			row.authorID,
			row.text,
			pgTimestamptz(row.createdAt),
			pgTimestamptz(row.updatedAt),
		)
	}

	return result

}

func pgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
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
