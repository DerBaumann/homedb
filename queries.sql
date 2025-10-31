-- name: CreateUser :one
INSERT INTO users (username, password, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users
WHERE username = $1 OR email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT id, username, email FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: ListItems :many
SELECT * FROM items
WHERE user_id = $1
ORDER BY name ASC;

-- name: CreateItem :one
INSERT INTO items (name, amount, unit, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
