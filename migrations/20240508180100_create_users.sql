-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL,
  email citext NOT NULL,
  encrypted_password varchar(500) NOT NULL,
  admin boolean NOT NULL DEFAULT false,
  display_name varchar(200) NOT NULL
);

CREATE UNIQUE INDEX idx_users_on_email ON users(email);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- +goose StatementEnd