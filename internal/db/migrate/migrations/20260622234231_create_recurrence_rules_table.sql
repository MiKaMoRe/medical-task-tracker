-- +goose Up
-- +goose StatementBegin
CREATE TYPE recurrence_rule_type AS ENUM (
  'weekly',
  'monthly',
  'yearly',
  'biweekly',
  'shift',
  'parity'
);

CREATE TABLE recurrence_rules (
  id        BIGSERIAL PRIMARY KEY,
  task_id   BIGINT NOT NULL REFERENCES tasks (id) ON DELETE CASCADE UNIQUE,
  rule_type recurrence_rule_type NOT NULL,
  end_date  DATE
);
CREATE INDEX idx_recurrence_rules_task_id ON recurrence_rules (task_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_recurrence_rules_task_id;
DROP TABLE IF EXISTS recurrence_rules;
DROP TYPE IF EXISTS recurrence_rule_type;
-- +goose StatementEnd
