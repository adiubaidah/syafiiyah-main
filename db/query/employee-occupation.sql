
-- name: QueryEmployeeOccupations :many
SELECT 
    *,
    COUNT(*) OVER () AS "count"
FROM 
    "employee_occupation";

-- name: CreateEmployeeOccupation :one
INSERT INTO "employee_occupation" ("name", "description") VALUES ($1, $2) RETURNING *;

-- name: UpdateEmployeeOccupation :one
UPDATE "employee_occupation" SET "name" = $1, "description" = $2 WHERE "id" = $3 RETURNING *;

-- name: DeleteEmployeeOccupation :one
DELETE FROM "employee_occupation" WHERE "id" = $1 RETURNING *;