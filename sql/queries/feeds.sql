-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *; 

-- name: GetAllFeeds :many
SELECT feeds.name AS feed_name, feeds.url AS feed_url, users.name AS user_name
FROM feeds JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY feeds.last_fetched_at NULLS FIRST
LIMIT 1;

-- name: MarkFeedFetched :exec
-- Update rows in 'feeds' where condition is met
UPDATE feeds
SET updated_at = $2, last_fetched_at = $2
WHERE feeds.id = $1;