-- name: CreateEmployeePresence :one
INSERT INTO
    "employee_presence" (
        "schedule_id",
        "schedule_name",
        "type",
        "employee_id",
        "notes",
        "created_by",
        "employee_permission_id"
    )
VALUES
    (
        @schedule_id,
        @schedule_name,
        @type :: presence_type,
        @employee_id,
        @notes,
        @created_by :: presence_created_by_type,
        @employee_permission_id
    ) RETURNING *;

-- name: CreateEmployeePresences :copyfrom
INSERT INTO
    "employee_presence" (
        "schedule_id",
        "schedule_name",
        "type",
        "employee_id",
        "notes",
        "created_at",
        "created_by",
        "employee_permission_id"
    )
VALUES
    (
        @schedule_id,
        @schedule_name,
        @type::presence_type,
        @employee_id,
        @notes,
        @created_at,
        @created_by :: presence_created_by_type,
        @employee_permission_id
    );

-- name: ListEmployeePresences :many
SELECT
    "employee_presence".*,
    "employee"."name" AS "employee_name"
FROM
    "employee_presence"
    INNER JOIN "employee" ON "employee_presence"."employee_id" = "employee"."id"
WHERE
    (
        sqlc.narg(employee_id) :: integer IS NULL
        OR "employee_id" = sqlc.narg(employee_id) :: integer
    )
    AND (
        sqlc.narg(q) :: text IS NULL
        OR "employee"."name" ILIKE '%' || sqlc.narg(q) || '%'
    )
    AND (
        sqlc.narg(type) :: presence_type IS NULL
        OR "type" = sqlc.narg(type) :: presence_type
    )
    AND (
        sqlc.narg(schedule_id) :: integer IS NULL
        OR "schedule_id" = sqlc.narg(schedule_id) :: integer
    )
    AND (
        sqlc.narg(from_date) :: date IS NULL
        OR DATE("created_at") >= sqlc.narg(from_date) :: date
    )
    AND (
        sqlc.narg(to_date) :: date IS NULL
        OR DATE("created_at") <= sqlc.narg(to_date) :: date
    )
ORDER BY
    "employee_presence"."id" DESC
LIMIT
    @limit_number OFFSET @offset_number;


-- name: CountEmployeePresences :one
SELECT
    COUNT(*)
FROM
    "employee_presence"
    INNER JOIN "employee" ON "employee_presence"."employee_id" = "employee"."id"
WHERE
    (
        sqlc.narg(employee_id) :: integer IS NULL
        OR "employee_id" = sqlc.narg(employee_id) :: integer
    )
    AND (
        sqlc.narg(q) :: text IS NULL
        OR "employee"."name" ILIKE '%' || sqlc.narg(q) || '%'
    )
    AND (
        sqlc.narg(type) :: presence_type IS NULL
        OR "type" = sqlc.narg(type) :: presence_type
    )
    AND (
        sqlc.narg(schedule_id) :: integer IS NULL
        OR "schedule_id" = sqlc.narg(schedule_id) :: integer
    )
    AND (
        sqlc.narg(from_date) :: date IS NULL
        OR DATE("created_at") >= sqlc.narg(from_date) :: date
    )
    AND (
        sqlc.narg(to_date) :: date IS NULL
        OR DATE("created_at") <= sqlc.narg(to_date) :: date
    );

-- name: ListMissingEmployeePresences :many
SELECT 
    "employee"."id", "employee"."name"
FROM
    "employee"
WHERE
    NOT EXISTS (
        SELECT
            1
        FROM
            "employee_presence"
        WHERE
            "employee_presence"."employee_id" = "employee"."id"
            AND DATE("employee_presence"."created_at") = sqlc.narg(date)::date
            AND "employee_presence"."schedule_id" = sqlc.narg(schedule_id)::integer
    );

-- name: UpdateEmployeePresence :one
UPDATE
    "employee_presence"
SET
    "schedule_id" = COALESCE(sqlc.narg(schedule_id), schedule_id),
    "schedule_name" = COALESCE(sqlc.narg(schedule_name), schedule_name),
    "type" = COALESCE(sqlc.narg(type)::presence_type, type),
    "employee_id" = COALESCE(sqlc.narg(employee_id), employee_id),
    "notes" = sqlc.narg(notes),
    "employee_permission_id" = @employee_permission_id
WHERE
    "id" = @id
RETURNING *;

-- name: DeleteEmployeePresence :one
DELETE FROM
    "employee_presence"
WHERE
    "id" = @id
RETURNING *;