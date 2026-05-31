-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts.comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts.posts(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES posts.comments(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL,
    text VARCHAR(2000) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_post_parent_created_at
    ON posts.comments (post_id, parent_id, created_at ASC);

CREATE OR REPLACE TRIGGER trigger_update_comment_timestamp
    BEFORE UPDATE ON posts.comments
    FOR EACH ROW
    EXECUTE FUNCTION posts.update_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_comment_timestamp ON posts.comments;
DROP INDEX IF EXISTS posts.idx_comments_post_parent_created_at;
DROP TABLE IF EXISTS posts.comments;
-- +goose StatementEnd
