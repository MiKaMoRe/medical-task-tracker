-- +goose Up
-- +goose StatementBegin
CREATE TABLE recurrence_monthly_rules (
  id                 BIGSERIAL PRIMARY KEY,
  recurrence_rule_id BIGINT NOT NULL UNIQUE REFERENCES recurrence_rules (id) ON DELETE CASCADE,
  month_day          SMALLINT NOT NULL CHECK (month_day >= 1 AND month_day <= 31)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS recurrence_monthly_rules;
