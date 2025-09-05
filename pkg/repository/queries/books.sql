-- name: GetBooks :many
select
  books.*,
  users.display_name as owner_display_name
from books
inner join users on users.id = books.owner_id
order by books.name, books.id
;

-- name: GetBookByID :one
select
  books.*,
  users.display_name as owner_display_name
from books
inner join users on users.id = books.owner_id
where books.id = $1
;

-- name: CreateBook :one
insert into books(
  created_at,
  updated_at,
  name,
  owner_id,
  default_currency_iso_code
) values (
  $1, -- created_at,
  $2, -- updated_at,
  $3, -- name,
  $4, -- owner_id,
  $5  -- default_currency_iso_code
)
returning *
;
