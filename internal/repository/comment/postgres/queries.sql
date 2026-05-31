-- name: AddComment :one
INSERT INTO posts.comments (
    post_id,
    author_id,
    parent_id,
    text
)
VALUES (sqlc.arg(post_id), sqlc.arg(author_id),
        sqlc.arg(parent_id), sqlc.arg(text))
RETURNING id, post_id, parent_id, author_id, text, created_at, updated_at;

-- name: GetComment :one
SELECT id, post_id, parent_id, author_id, text, created_at, updated_at
FROM posts.comments
WHERE id = sqlc.arg(id);

-- name: GetComments :many
SELECT id, post_id, parent_id, author_id, text, created_at, updated_at
FROM posts.comments
WHERE post_id = sqlc.arg(post_id)
  AND parent_id IS NOT DISTINCT FROM sqlc.arg(parent_id)
ORDER BY created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
