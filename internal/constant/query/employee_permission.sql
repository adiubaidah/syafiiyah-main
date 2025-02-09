-- name: CreateEmployeePermission :one
INSERT INTO
    "employee_permission" (
        employee_id,
        start_permission,
        end_permission,
        "type",
        excuse
    )
VALUES
    (
        @employee_id,
        @start_permission,
        sqlc.narg(end_permission),
        @type :: santri_permission_type,
        @excuse
    ) RETURNING *;

-- name: ListEmployeePermissions :many
SELECT
    "employee_permission".*,
    "employee"."name" AS "employee_name"
FROM
    "employee_permission"
    INNER JOIN "employee" ON "employee_permission"."employee_id" = "employee"."id"
WHERE
    (sqlc.narg(q) :: text IS NULL
    OR "employee"."name" ILIKE '%' || sqlc.narg(q) || '%')
    AND (
        sqlc.narg(employee_id) :: integer IS NULL
        OR "employee_id" = sqlc.narg(employee_id) :: integer
    )
    AND (
        sqlc.narg(type) :: santri_permission_type IS NULL
        OR "type" = sqlc.narg(type) :: santri_permission_type
    )
    AND (
        sqlc.narg(from_date) :: timestamptz IS NULL
        OR "start_permission" >= sqlc.narg(from_date) :: timestamptz
    )
    AND (
        sqlc.narg(end_date) :: timestamptz IS NULL
        OR "end_permission" <= sqlc.narg(end_date) :: timestamptz
    )
LIMIT
    @limit_number OFFSET @offset_number;

-- name: GetEmployeePermission :one
SELECT
    "employee_permission".*,
    "employee"."name" AS "employee_name"
FROM
    "employee_permission"
    INNER JOIN "employee" ON "employee_permission"."employee_id" = "employee"."id"
WHERE
    "employee_permission"."id" = @id;

-- name: UpdateEmployeePermission :one
UPDATE
    "employee_permission"
SET
    "employee_id" = COALESCE(sqlc.narg(employee_id), employee_id),
    "start_permission" = COALESCE(sqlc.narg(start_permission), start_permission),
    "end_permission" = sqlc.narg(end_permission),
    "excuse" = COALESCE(sqlc.narg(excuse), excuse)
WHERE
    "id" = @id RETURNING *;

-- name: DeleteEmployeePermission :one
DELETE FROM
    "employee_permission"
WHERE
    "id" = @id RETURNING *;