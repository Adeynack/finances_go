-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS "citext";

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "uuid-ossp";

DROP EXTENSION IF EXISTS "citext";

-- +goose StatementEnd