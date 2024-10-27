
-- name: ListEmployeeOccupations :many
SELECT 
    *,
    COUNT(*) OVER () AS "count"
FROM 
    "employee_occupation";

-- name: CreateEmployeeOccupation :one
INSERT INTO "employee_occupation" ("name", "description") VALUES (@name, @description) RETURNING *;

-- name: UpdateEmployeeOccupation :one
UPDATE "employee_occupation" SET "name" = @name, "description" = @description WHERE "id" = @id RETURNING *;

-- name: DeleteEmployeeOccupation :one
DELETE FROM "employee_occupation" WHERE "id" = @id RETURNING *;