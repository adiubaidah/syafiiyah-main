-- name: CreateHoliday :one
INSERT INTO "holiday" ("name", "color", "description") VALUES (@name, @color, @description) RETURNING *;

-- name: UpdateHoliday :one
UPDATE "holiday" SET "name" = @name, "color" = @color, "description" = @description WHERE "id" = @id RETURNING *;

-- name: DeleteHoliday :one
DELETE FROM "holiday" WHERE "id" = @id RETURNING *;
