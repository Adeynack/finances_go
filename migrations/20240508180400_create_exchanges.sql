-- +goose Up
-- +goose StatementBegin
CREATE TYPE exchange_status AS ENUM ('uncleared', 'reconciling', 'cleared');

CREATE TABLE exchanges (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL,
  date date NOT NULL,
  register_id uuid NOT NULL REFERENCES registers(id),
  cheque varchar(100),
  description varchar(1_000) NOT NULL,
  memo varchar(1_000_000),
  status exchange_status NOT NULL DEFAULT 'uncleared'
);

CREATE INDEX ON exchanges(date);

COMMENT ON COLUMN exchanges.date IS 'Date the exchange appears in the book.';

COMMENT ON COLUMN exchanges.register_id IS 'From which register does the money come from.';

COMMENT ON COLUMN exchanges.cheque IS 'Cheque information.';

COMMENT ON COLUMN exchanges.description IS 'Label of the exchange.';

COMMENT ON COLUMN exchanges.memo IS 'Detail about the exchange.';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exchanges CASCADE;

DROP TYPE IF EXISTS exchange_status CASCADE;

-- +goose StatementEnd