-- +goose Up
-- +goose StatementBegin
CREATE TABLE recurrence_parity_rules (
  id                 BIGSERIAL PRIMARY KEY,
  recurrence_rule_id BIGINT NOT NULL UNIQUE REFERENCES recurrence_rules (id) ON DELETE CASCADE,
  is_even            BOOLEAN NOT NULL DEFAULT TRUE
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS recurrence_parity_rules;
