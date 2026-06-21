-- name: CreateUser :one
INSERT INTO users (
    name, email, password_hash, role
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, name, email, password_hash, role, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, name, email, password_hash, role, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, role, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
    name = COALESCE($2, name),
    email = COALESCE($3, email),
    password_hash = COALESCE($4, password_hash),
    role = COALESCE($5, role),
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, email, password_hash, role, created_at, updated_at;

-- name: ListUsers :many
SELECT id, name, email, password_hash, role, created_at, updated_at
FROM users
ORDER BY created_at DESC;
