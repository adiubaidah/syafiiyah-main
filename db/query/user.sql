-- name: CreateUser :one
INSERT INTO "user" ("role", "username", "password") VALUES ($1, $2, $3) RETURNING *;

-- -- name: QueryUser :many
-- SELECT "user"."id", "user"."username", "user"."role" FROM "user"
-- LEFT JOIN "parent" ON "user"."id" = "parent"."user_id"
-- LEFT JOIN "employee" ON "user"."id" = "employee"."user_id"
-- WHERE "user"."role" = $1