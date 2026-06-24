-- +goose Up
-- +goose StatementBegin
CREATE TABLE tags_tasks (
  task_id BIGINT NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
  tag_id  BIGINT NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
  PRIMARY KEY (task_id, tag_id)
);

CREATE INDEX idx_tags_tasks_task_id ON tags_tasks (task_id);
CREATE INDEX idx_tags_tasks_tag_id ON tags_tasks (tag_id);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS tags_tasks;
