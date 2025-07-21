-- +goose Up
-- +goose StatementBegin
INSERT INTO short_urls (short_path, original_url) VALUES
('/urlshort', 'https://github.com/gophercises/urlshort'),
('/urlshort-final', 'https://github.com/gophercises/urlshort/tree/solution');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM short_urls WHERE short_path IN ('/urlshort', '/urlshort-final');
-- +goose StatementEnd
