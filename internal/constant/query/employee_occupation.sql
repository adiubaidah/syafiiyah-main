-- name: ListEmployeeOccupations :many
SELECT
    "employee_occupation".*,
    COUNT("employee"."id") AS "count"
FROM
    "employee_occupation"
    LEFT JOIN "employee" ON "employee"."occupation_id" = "employee_occupation"."id"
GROUP BY
    "employee_occupation"."id"
ORDER BY
    "employee_occupation"."id" ASC;

-- name: CreateEmployeeOccupation :one
INSERT INTO
    "employee_occupation" ("name", "description")
VALUES
    (@name, @description) RETURNING *;

-- name: UpdateEmployeeOccupation :one
UPDATE
    "employee_occupation"
SET
    "name" = @name,
    "description" = @description
WHERE
    "id" = @id RETURNING *;

-- name: DeleteEmployeeOccupation :one
DELETE FROM
    "employee_occupation"
WHERE
    "id" = @id RETURNING *;