-- +goose Up
-- +goose StatementBegin
CREATE TABLE splits (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL,
  exchange_id uuid NOT NULL REFERENCES exchanges(id),
  destination_register_id uuid NOT NULL REFERENCES registers(id),
  amount bigint NOT NULL DEFAULT 0,
  counterpart_amount bigint,
  memo varchar(1_000_000),
  status exchange_status NOT NULL DEFAULT 'uncleared'
);

COMMENT ON COLUMN splits.destination_register_id IS 'To which register is the money going to for this split.';

COMMENT ON COLUMN splits.counterpart_amount IS 'Amount in the destination register, if it differs from ''amount'' (ex: an exchange rate applies).';

COMMENT ON COLUMN splits.memo IS 'Detail about the exchange, to show in the destination register.';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS splits CASCADE;

-- +goose StatementEnd