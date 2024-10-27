-- name: CreateSantriOccupation :one
INSERT INTO "santri_occupation" ("name", "description") VALUES (@name, @description) RETURNING *;

-- name: QuerySantriOccupations :many
SELECT 
    *,
    COUNT(*) OVER () AS "count"
FROM
    "santri_occupation";


-- name: UpdateSantriOccupation :one
UPDATE "santri_occupation" SET "name" = $1, "description" = $2 WHERE "id" = $3 RETURNING *;

-- name: DeleteSantriOccupation :one
DELETE FROM "santri_occupation" WHERE "id" = $1 RETURNING *;