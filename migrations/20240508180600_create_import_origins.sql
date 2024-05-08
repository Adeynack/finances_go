-- +goose Up
-- +goose StatementBegin
CREATE TABLE import_origins (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL,
  subject_type varchar NOT NULL,
  subject_id uuid NOT NULL,
  external_system varchar NOT NULL,
  external_id varchar NOT NULL
);

CREATE INDEX ON import_origins(subject_type, subject_id);

CREATE UNIQUE INDEX ON import_origins(subject_type, subject_id, external_system);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS import_origins CASCADE;

-- +goose StatementEnd