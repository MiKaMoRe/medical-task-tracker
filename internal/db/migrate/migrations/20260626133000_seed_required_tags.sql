-- +goose Up
-- +goose StatementBegin
INSERT INTO tags (name, is_required)
VALUES
  ('отчётность', TRUE),
  ('операции', TRUE),
  ('звонок', TRUE)
ON CONFLICT (name) DO UPDATE
SET is_required = EXCLUDED.is_required;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tags
WHERE name IN ('отчётность', 'операции', 'звонок')
  AND is_required = TRUE;
-- +goose StatementEnd
