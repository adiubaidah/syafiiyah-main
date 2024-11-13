-- name: CreateSantriOccupation :one
INSERT INTO
    "santri_occupation" ("name", "description")
VALUES
    (@name, @description) RETURNING *;

-- name: ListSantriOccupations :many
SELECT
    "santri_occupation".*,
    COUNT("santri"."id") AS "count"
FROM
    "santri_occupation"
    LEFT JOIN "santri" ON "santri"."occupation_id" = "santri_occupation"."id"
GROUP BY
    "santri_occupation"."id"
ORDER BY
    "santri_occupation"."id" ASC;

-- name: UpdateSantriOccupation :one
UPDATE
    "santri_occupation"
SET
    "name" = COALESCE(sqlc.narg(name), name),
    "description" = sqlc.narg(description)
WHERE
    "id" = @id RETURNING *;

-- name: DeleteSantriOccupation :one
DELETE FROM
    "santri_occupation"
WHERE
    "id" = @id RETURNING *;