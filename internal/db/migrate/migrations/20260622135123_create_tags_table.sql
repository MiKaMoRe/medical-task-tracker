-- +goose Up
-- +goose StatementBegin
CREATE TABLE tags (
  id          BIGSERIAL PRIMARY KEY,
  name        VARCHAR(255) NOT NULL UNIQUE,
  is_required BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS tags;
