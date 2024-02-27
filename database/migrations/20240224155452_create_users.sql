-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "email" varchar NOT NULL,
  "encrypted_password" varchar NOT NULL,
  "admin" bool NOT NULL DEFAULT false,
  "display_name" varchar NOT NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";

-- +goose StatementEnd