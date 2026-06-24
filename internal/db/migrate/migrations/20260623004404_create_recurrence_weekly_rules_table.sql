-- +goose Up
-- +goose StatementBegin
CREATE TABLE recurrence_weekly_rules (
  id                 BIGSERIAL PRIMARY KEY,
  recurrence_rule_id BIGINT NOT NULL UNIQUE REFERENCES recurrence_rules (id) ON DELETE CASCADE,
  week_day           SMALLINT NOT NULL CHECK (week_day >= 0 AND week_day < 7)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS recurrence_weekly_rules;
