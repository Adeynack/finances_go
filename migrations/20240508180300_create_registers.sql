-- +goose Up
-- +goose StatementBegin
CREATE TYPE register_types AS ENUM (
  'Bank',
  'Card',
  'Investment',
  'Asset',
  'Liability',
  'Loan',
  'Institution',
  'Expense',
  'Income'
);

CREATE TABLE registers (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL,
  name varchar(500) NOT NULL,
  type register_types NOT NULL,
  book_id uuid NOT NULL REFERENCES books(id),
  parent_id uuid REFERENCES registers(id),
  starts_at date NOT NULL,
  expires_at date,
  currency_iso_code char(3) NOT NULL,
  notes varchar(1_000_000),
  initial_balance bigint NOT NULL DEFAULT 0,
  active boolean NOT NULL DEFAULT true,
  default_category uuid REFERENCES registers(id),
  institution_name varchar(200),
  account_number varchar(500),
  iban varchar(34),
  annual_interest_rate numeric,
  credit_limit bigint,
  card_number varchar(200)
);

COMMENT ON COLUMN registers.parent_id IS 'A null parent means it is a root register.';

COMMENT ON COLUMN registers.starts_at IS 'Opening date of the register.';

COMMENT ON COLUMN registers.expires_at IS 'Optional expiration date of the register (ex: for a credit card).';

COMMENT ON COLUMN registers.initial_balance IS 'Balance amount in the account at its creation (typically 0)';

COMMENT ON COLUMN registers.default_category IS 'The category automatically selected when entering a new exchange from this register.';

COMMENT ON COLUMN registers.institution_name IS 'Name of the institution (ex: bank) managing the registry (ex: credit card).';

COMMENT ON COLUMN registers.account_number IS 'account_numberNumber by which the register is referred to (ex: bank account number).';

COMMENT ON COLUMN registers.iban IS 'In the case the register is identified by an International Bank Account Number (IBAN).';

COMMENT ON COLUMN registers.annual_interest_rate IS 'In the case the register is being charged interests, its rate per year (ex: credit card).';

COMMENT ON COLUMN registers.credit_limit IS 'In the case the register has a credit limit (ex: credit card, credit margin).';

COMMENT ON COLUMN registers.card_number IS 'In the case the register is linked to a card, its number (ex: a credit card).';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS registers CASCADE;

DROP TYPE IF EXISTS register_types CASCADE;

-- +goose StatementEnd