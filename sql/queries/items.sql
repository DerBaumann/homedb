-- name: ListItems :many
SELECT * FROM items
WHERE user_id = $1
ORDER BY name ASC;

-- name: FilterItemsByName :many
SELECT * FROM items
WHERE user_id = $1
AND name ILIKE $2
ORDER BY name ASC;

-- name: GetItemByID :one
SELECT * FROM items
WHERE id = $1
LIMIT 1;

-- name: CreateItem :one
INSERT INTO items (name, amount, unit, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateItem :one
UPDATE items
SET
  name = $1,
  amount = $2,
  unit = $3
WHERE id = $4
RETURNING *;

-- name: DeleteItem :one
DELETE FROM items
WHERE id = $1
RETURNING *;
