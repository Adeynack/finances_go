-- name: GetExchanges :many
select
  exchanges.*
from
  exchanges
order by
  exchanges.date,
  exchanges.created_at,
  exchanges.id
;
