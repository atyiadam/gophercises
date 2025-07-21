-- name: GetShortURLByPath :one
SELECT id, short_path, original_url, created_at FROM short_urls
WHERE short_path = $1 LIMIT 1;
