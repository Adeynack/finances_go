-- +goose Up
-- +goose StatementBegin
CREATE TABLE books (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL,
  name varchar(1000) NOT NULL,
  owner_id uuid NOT NULL REFERENCES users(id),
  default_currency_iso_code char(3) NOT NULL
);

CREATE INDEX idx_books_on_owner_id ON books(owner_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS books;

-- +goose StatementEnd