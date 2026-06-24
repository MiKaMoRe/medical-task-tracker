-- +goose Up
-- +goose StatementBegin
CREATE TABLE task_occurrence_completions (
  id              BIGSERIAL NOT NULL UNIQUE,
  task_id         BIGINT NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
  occurrence_date DATE NOT NULL,
  completed_at    TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (task_id, occurrence_date)
);
CREATE INDEX idx_task_occurrence_completions_task_id ON task_occurrence_completions (task_id);
-- +goose StatementEnd

-- +goose Down
DROP INDEX IF EXISTS idx_task_occurrence_completions_task_id;
DROP TABLE IF EXISTS task_occurrence_completions;
