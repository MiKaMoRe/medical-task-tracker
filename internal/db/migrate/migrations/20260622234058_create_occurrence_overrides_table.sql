-- +goose Up
-- +goose StatementBegin
CREATE TYPE occurrence_type AS ENUM ('skip', 'add', 'reschedule');

CREATE TABLE occurrence_overrides (
  task_id         BIGINT NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
  occurrence_type occurrence_type NOT NULL,
  occurrence_date DATE NOT NULL,
  reschedule_at   DATE,
  PRIMARY KEY (task_id, occurrence_date),

  CONSTRAINT occurrence_overrides_reschedule_at_check CHECK (
    (occurrence_type = 'reschedule' AND reschedule_at IS NOT NULL)
    OR (occurrence_type != 'reschedule' AND reschedule_at IS NULL)
  )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS occurrence_overrides;
DROP TYPE IF EXISTS occurrence_type;
-- +goose StatementEnd
