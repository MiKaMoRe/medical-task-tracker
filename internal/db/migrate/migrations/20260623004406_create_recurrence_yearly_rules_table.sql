-- +goose Up
-- +goose StatementBegin
CREATE TABLE recurrence_yearly_rules (
  id                 BIGSERIAL PRIMARY KEY,
  recurrence_rule_id BIGINT NOT NULL UNIQUE REFERENCES recurrence_rules (id) ON DELETE CASCADE,
  month              SMALLINT NOT NULL CHECK (month >= 1 AND month <= 12),
  day                SMALLINT NOT NULL CHECK (day >= 1 AND day <= 31)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS recurrence_yearly_rules;
