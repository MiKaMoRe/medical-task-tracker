-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
  id            BIGSERIAL PRIMARY KEY,
  name          VARCHAR(255) NOT NULL,
  description   TEXT,
  date          TIMESTAMPTZ NOT NULL,
  is_recurrence BOOLEAN NOT NULL DEFAULT FALSE,
  is_done       BOOLEAN NOT NULL DEFAULT FALSE,
  completed_at  TIMESTAMPTZ,

  CONSTRAINT tasks_recurring_not_marked_done
    CHECK (NOT is_recurrence OR is_done = FALSE)
);

CREATE INDEX idx_tasks_date ON tasks (date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tasks_date;
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
