-- name: GetUserByID :one
select
  users.*
from
  users
where
  users.id = $1
;
