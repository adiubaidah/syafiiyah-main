-- name: CreateArduino :one
INSERT INTO "arduino" ("name") VALUES (@name) RETURNING *;

-- name: QueryArduino :many
SELECT * FROM "arduino" WHERE "name" ILIKE '%' || @name || '%' LIMIT @limit_number OFFSET @offset_number;

-- name: UpdateArduino :one
UPDATE "arduino" SET "name" = @name WHERE "id" = @id RETURNING *;

-- name: DeleteArduino :one
DELETE FROM "arduino" WHERE "id" = @id RETURNING *;