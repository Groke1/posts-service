-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS posts;

CREATE TABLE posts.posts (
     id BIGSERIAL PRIMARY KEY,
     title VARCHAR(500) NOT NULL,
     text TEXT NOT NULL,
     author_id BIGINT NOT NULL,
     comments_enabled BOOLEAN NOT NULL DEFAULT TRUE,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_author_created_at
    ON posts.posts USING BTREE (author_id, created_at DESC);

CREATE INDEX idx_posts_created_at
    ON posts.posts USING BTREE (created_at DESC);

CREATE OR REPLACE FUNCTION posts.update_timestamp()
RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_posts_timestamp
    BEFORE UPDATE ON posts.posts
    FOR EACH ROW
    EXECUTE FUNCTION posts.update_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_posts_timestamp ON posts.posts;
DROP FUNCTION IF EXISTS posts.update_timestamp;
DROP INDEX IF EXISTS posts.idx_posts_created_at;
DROP INDEX IF EXISTS posts.idx_posts_author_created_at;
DROP TABLE IF EXISTS posts.posts;
DROP SCHEMA IF EXISTS posts;
-- +goose StatementEnd
