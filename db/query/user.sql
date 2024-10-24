-- name: CreateUser :one
INSERT INTO "user" ("role", "username", "password") VALUES ($1, $2, $3) RETURNING *;

-- name: GetUserByUsernameAndRole :many
SELECT * FROM "user"
WHERE "username" = $1
AND ($2::text IS NULL OR "role" = $2)
ORDER BY "username"
LIMIT $3 OFFSET $4;