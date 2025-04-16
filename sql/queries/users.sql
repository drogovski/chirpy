-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES(
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserCredentials :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpgradeToChirpyRed :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1;
