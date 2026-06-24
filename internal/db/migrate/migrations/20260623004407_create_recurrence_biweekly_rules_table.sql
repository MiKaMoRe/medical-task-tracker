-- +goose Up
-- +goose StatementBegin
CREATE TABLE recurrence_biweekly_rules (
  id                 BIGSERIAL PRIMARY KEY,
  recurrence_rule_id BIGINT NOT NULL UNIQUE REFERENCES recurrence_rules (id) ON DELETE CASCADE,
  is_odd             BOOLEAN NOT NULL,
  week_day           SMALLINT NOT NULL CHECK (week_day >= 0 AND week_day < 7)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS recurrence_biweekly_rules;
