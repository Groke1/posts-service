-- name: AddPost :one
INSERT INTO posts.posts (
    author_id,
    title,
    text,
    comments_enabled
)
VALUES (sqlc.arg(author_id), sqlc.arg(title),
        sqlc.arg(text), sqlc.arg(comments_enabled))
RETURNING id, author_id, title, text, comments_enabled, created_at, updated_at;

-- name: GetPost :one
SELECT id, author_id, title, text, comments_enabled, created_at, updated_at
FROM posts.posts
WHERE id = sqlc.arg(id);

-- name: GetPostForUpdate :one
SELECT id, author_id, title, text, comments_enabled, created_at, updated_at
FROM posts.posts
WHERE id = sqlc.arg(id)
FOR UPDATE;

-- name: GetPostsByAuthor :many
SELECT id, author_id, title, text, comments_enabled, created_at, updated_at
FROM posts.posts
WHERE author_id = sqlc.arg(author_id)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetPosts :many
SELECT id, author_id, title, text, comments_enabled, created_at, updated_at
FROM posts.posts
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: SetCommentsEnabled :one
UPDATE posts.posts
SET comments_enabled = sqlc.arg(comments_enabled)
WHERE id = sqlc.arg(id)
RETURNING id, author_id, title, text, comments_enabled, created_at, updated_at;
