-- +goose Up
-- +goose StatementBegin
CREATE TABLE recurrence_shift_rules (
  id                   BIGSERIAL PRIMARY KEY,
  recurrence_rule_id   BIGINT NOT NULL UNIQUE REFERENCES recurrence_rules (id) ON DELETE CASCADE,
  number_of_task_days  INT NOT NULL,
  number_of_shift_days INT NOT NULL,

  CONSTRAINT chk_recurrence_shift_rules_task_days
    CHECK (number_of_task_days > 0),
  CONSTRAINT chk_recurrence_shift_rules_shift_days
    CHECK (number_of_shift_days >= 0)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS recurrence_shift_rules;
