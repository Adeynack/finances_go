-- name: GetSplitsForExchangeIds :many
select
  splits.amount,
  splits.counterpart_amount,
  splits.created_at,
  splits.destination_register_id,
  splits.exchange_id,
  splits.id,
  splits.memo,
  splits.status,
  splits.updated_at
from
  splits
where
  splits.exchange_id = ANY($1::uuid[])
order by
  splits.created_at,
  splits.id
;
