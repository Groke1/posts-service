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

func TestPostRepository_GetPostByID(t *testing.T) {
	t.Parallel()

	now := time.Now()

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
				expectGetPost(deps.dbMock, 100, postRow{
					id:              100,
					authorID:        42,
					title:           "Post title",
					text:            "Post text",
					commentsEnabled: true,
					createdAt:       now,
					updatedAt:       now,
				}, nil)
			},
			expectedPost: &entity.Post{
				ID:              100,
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: true,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			expectedErr: nil,
		},
		{
			name: "post not found",
			id:   100,
			setupMocks: func(deps testDeps) {
				expectGetPost(deps.dbMock, 100, postRow{}, pgx.ErrNoRows)
			},
			expectedPost: nil,
			expectedErr:  entity.ErrPostNotFound,
		},
		{
			name: "db error",
			id:   100,
			setupMocks: func(deps testDeps) {
				expectGetPost(deps.dbMock, 100, postRow{}, errDB)
			},
			expectedPost: nil,
			expectedErr:  errDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			post, err := deps.repo.GetPostByID(context.Background(), tc.id)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedPost, post)
		})
	}
}

func TestPostRepository_GetPostByIDForUpdate(t *testing.T) {
	t.Parallel()

	now := time.Now()

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
				expectGetPostForUpdate(deps.dbMock, 100, postRow{
					id:              100,
					authorID:        42,
					title:           "Post title",
					text:            "Post text",
					commentsEnabled: true,
					createdAt:       now,
					updatedAt:       now,
				}, nil)
			},
			expectedPost: &entity.Post{
				ID:              100,
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: true,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			expectedErr: nil,
		},
		{
			name: "post not found",
			id:   100,
			setupMocks: func(deps testDeps) {
				expectGetPostForUpdate(deps.dbMock, 100, postRow{}, pgx.ErrNoRows)
			},
			expectedPost: nil,
			expectedErr:  entity.ErrPostNotFound,
		},
		{
			name: "db error",
			id:   100,
			setupMocks: func(deps testDeps) {
				expectGetPostForUpdate(deps.dbMock, 100, postRow{}, errDB)
			},
			expectedPost: nil,
			expectedErr:  errDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			post, err := deps.repo.GetPostByIDForUpdate(context.Background(), tc.id)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedPost, post)
		})
	}
}

func TestPostRepository_GetPostsByAuthorID(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testCases := []struct {
		name          string
		authorID      int64
		limit         int
		offset        int
		setupMocks    mocksSetupFunc
		expectedPosts []*entity.Post
		expectedErr   error
	}{
		{
			name:     "success",
			authorID: 42,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPostsByAuthor(deps.dbMock, 42, 10, 0, []postRow{
					{
						id:              100,
						authorID:        42,
						title:           "Post title 1",
						text:            "Post text 1",
						commentsEnabled: true,
						createdAt:       now,
						updatedAt:       now,
					},
					{
						id:              101,
						authorID:        42,
						title:           "Post title 2",
						text:            "Post text 2",
						commentsEnabled: false,
						createdAt:       now,
						updatedAt:       now,
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
					CreatedAt:       now,
					UpdatedAt:       now,
				},
				{
					ID:              101,
					AuthorID:        42,
					Title:           "Post title 2",
					Text:            "Post text 2",
					CommentsEnabled: false,
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			expectedErr: nil,
		},
		{
			name:     "empty result",
			authorID: 42,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPostsByAuthor(deps.dbMock, 42, 10, 0, nil, nil)
			},
			expectedPosts: []*entity.Post{},
			expectedErr:   nil,
		},
		{
			name:     "post not found",
			authorID: 42,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPostsByAuthor(deps.dbMock, 42, 10, 0, nil, pgx.ErrNoRows)
			},
			expectedPosts: nil,
			expectedErr:   entity.ErrPostNotFound,
		},
		{
			name:     "db error",
			authorID: 42,
			limit:    10,
			offset:   0,
			setupMocks: func(deps testDeps) {
				expectGetPostsByAuthor(deps.dbMock, 42, 10, 0, nil, errDB)
			},
			expectedPosts: nil,
			expectedErr:   errDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			posts, err := deps.repo.GetPostsByAuthorID(
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

func TestPostRepository_GetPosts(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testCases := []struct {
		name          string
		limit         int
		offset        int
		setupMocks    mocksSetupFunc
		expectedPosts []*entity.Post
		expectedErr   error
	}{
		{
			name:   "success",
			limit:  10,
			offset: 0,
			setupMocks: func(deps testDeps) {
				expectGetPosts(deps.dbMock, 10, 0, []postRow{
					{
						id:              100,
						authorID:        42,
						title:           "Post title 1",
						text:            "Post text 1",
						commentsEnabled: true,
						createdAt:       now,
						updatedAt:       now,
					},
					{
						id:              101,
						authorID:        43,
						title:           "Post title 2",
						text:            "Post text 2",
						commentsEnabled: false,
						createdAt:       now,
						updatedAt:       now,
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
					CreatedAt:       now,
					UpdatedAt:       now,
				},
				{
					ID:              101,
					AuthorID:        43,
					Title:           "Post title 2",
					Text:            "Post text 2",
					CommentsEnabled: false,
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			expectedErr: nil,
		},
		{
			name:   "empty result",
			limit:  10,
			offset: 0,
			setupMocks: func(deps testDeps) {
				expectGetPosts(deps.dbMock, 10, 0, nil, nil)
			},
			expectedPosts: []*entity.Post{},
			expectedErr:   nil,
		},
		{
			name:   "db error",
			limit:  10,
			offset: 0,
			setupMocks: func(deps testDeps) {
				expectGetPosts(deps.dbMock, 10, 0, nil, errDB)
			},
			expectedPosts: nil,
			expectedErr:   errDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			posts, err := deps.repo.GetPosts(context.Background(), tc.limit, tc.offset)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedPosts, posts)
		})
	}
}

func TestPostRepository_CreatePost(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testCases := []struct {
		name            string
		authorID        int64
		title           string
		text            string
		commentsEnabled bool
		setupMocks      mocksSetupFunc
		expectedPost    *entity.Post
		expectedErr     error
	}{
		{
			name:            "success",
			authorID:        42,
			title:           "Post title",
			text:            "Post text",
			commentsEnabled: true,
			setupMocks: func(deps testDeps) {
				expectAddPost(deps.dbMock, 42, "Post title", "Post text", true, postRow{
					id:              100,
					authorID:        42,
					title:           "Post title",
					text:            "Post text",
					commentsEnabled: true,
					createdAt:       now,
					updatedAt:       now,
				}, nil)
			},
			expectedPost: &entity.Post{
				ID:              100,
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: true,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			expectedErr: nil,
		},
		{
			name:            "db error",
			authorID:        42,
			title:           "Post title",
			text:            "Post text",
			commentsEnabled: true,
			setupMocks: func(deps testDeps) {
				expectAddPost(deps.dbMock, 42, "Post title", "Post text", true, postRow{}, errDB)
			},
			expectedPost: nil,
			expectedErr:  errDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			post, err := deps.repo.CreatePost(
				context.Background(),
				tc.authorID,
				tc.title,
				tc.text,
				tc.commentsEnabled,
			)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedPost, post)
		})
	}
}

func TestPostRepository_SetPostCommentsEnabled(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testCases := []struct {
		name         string
		id           int64
		enabled      bool
		setupMocks   mocksSetupFunc
		expectedPost *entity.Post
		expectedErr  error
	}{
		{
			name:    "success",
			id:      100,
			enabled: false,
			setupMocks: func(deps testDeps) {
				expectSetCommentsEnabled(deps.dbMock, 100, false, postRow{
					id:              100,
					authorID:        42,
					title:           "Post title",
					text:            "Post text",
					commentsEnabled: false,
					createdAt:       now,
					updatedAt:       now,
				}, nil)
			},
			expectedPost: &entity.Post{
				ID:              100,
				AuthorID:        42,
				Title:           "Post title",
				Text:            "Post text",
				CommentsEnabled: false,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			expectedErr: nil,
		},
		{
			name:    "db error",
			id:      100,
			enabled: false,
			setupMocks: func(deps testDeps) {
				expectSetCommentsEnabled(deps.dbMock, 100, false, postRow{}, errDB)
			},
			expectedPost: nil,
			expectedErr:  errDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tc.setupMocks(deps)

			post, err := deps.repo.SetPostCommentsEnabled(
				context.Background(),
				tc.id,
				tc.enabled,
			)

			checkError(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedPost, post)
		})
	}
}

type mocksSetupFunc func(deps testDeps)

type testDeps struct {
	dbMock pgxmock.PgxConnIface
	repo   *postRepository
}

type postRow struct {
	id              int64
	authorID        int64
	title           string
	text            string
	commentsEnabled bool
	createdAt       time.Time
	updatedAt       time.Time
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
		repo:   NewPostRepository(dbMock),
	}
}

func expectGetPost(
	dbMock pgxmock.PgxConnIface,
	id int64,
	row postRow,
	err error,
) {
	query := dbMock.ExpectQuery(".*").
		WithArgs(id)

	if err != nil {
		query.WillReturnError(err)
		return
	}

	query.WillReturnRows(newPostRows(row))
}

func expectGetPostForUpdate(
	dbMock pgxmock.PgxConnIface,
	id int64,
	row postRow,
	err error,
) {
	query := dbMock.ExpectQuery(".*").
		WithArgs(id)

	if err != nil {
		query.WillReturnError(err)
		return
	}

	query.WillReturnRows(newPostRows(row))
}

func expectGetPostsByAuthor(
	dbMock pgxmock.PgxConnIface,
	authorID int64,
	limit int,
	offset int,
	rows []postRow,
	err error,
) {
	query := dbMock.ExpectQuery(".*").
		WithArgs(authorID, int32(offset), int32(limit))

	if err != nil {
		query.WillReturnError(err)
		return
	}

	query.WillReturnRows(newPostsRows(rows))
}

func expectGetPosts(
	dbMock pgxmock.PgxConnIface,
	limit int,
	offset int,
	rows []postRow,
	err error,
) {
	query := dbMock.ExpectQuery(".*").
		WithArgs(int32(offset), int32(limit))

	if err != nil {
		query.WillReturnError(err)
		return
	}

	query.WillReturnRows(newPostsRows(rows))
}

func expectAddPost(
	dbMock pgxmock.PgxConnIface,
	authorID int64,
	title string,
	text string,
	commentsEnabled bool,
	row postRow,
	err error,
) {
	query := dbMock.ExpectQuery(".*").
		WithArgs(authorID, title, text, commentsEnabled)

	if err != nil {
		query.WillReturnError(err)
		return
	}

	query.WillReturnRows(newPostRows(row))
}

func expectSetCommentsEnabled(
	dbMock pgxmock.PgxConnIface,
	id int64,
	enabled bool,
	row postRow,
	err error,
) {
	query := dbMock.ExpectQuery(".*").
		WithArgs(enabled, id)

	if err != nil {
		query.WillReturnError(err)
		return
	}

	query.WillReturnRows(newPostRows(row))
}

func newPostRows(row postRow) *pgxmock.Rows {
	return newPostsRows([]postRow{row})
}

func newPostsRows(rows []postRow) *pgxmock.Rows {
	result := pgxmock.NewRows([]string{
		"id", "author_id",
		"title", "text",
		"comments_enabled",
		"created_at", "updated_at",
	})

	for _, row := range rows {
		result.AddRow(
			row.id, row.authorID,
			row.title, row.text,
			row.commentsEnabled,
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
