-- name: CreateRfid :one
INSERT INTO
    rfid ("uid", "is_active", "santri_id", "employee_id")
VALUES
    (
        @uid,
        @is_active,
        sqlc.narg(santri_id),
        sqlc.narg(employee_id)
    ) RETURNING *;

-- name: ListRfid :many
SELECT
    "rfid".*,
    "santri"."name" as "santri_name",
    "employee"."name" as "employee_name"
FROM
    rfid
    LEFT JOIN "santri" ON "rfid"."santri_id" = "santri"."id"
    LEFT JOIN "employee" ON "rfid"."employee_id" = "employee"."id"
WHERE
    (
        sqlc.narg(q) :: text IS NULL
        OR "uid" ILIKE '%' || sqlc.narg(q) || '%'
        OR "santri"."name" ILIKE '%' || sqlc.narg(q) || '%'
        OR "employee"."name" ILIKE '%' || sqlc.narg(q) || '%'
    )
    AND (
        sqlc.narg(is_active)::boolean IS NULL
        OR "rfid"."is_active" = sqlc.narg(is_active)
    )
    AND (
        (
            sqlc.narg(is_santri)::boolean IS NULL
            OR sqlc.narg(is_santri) = FALSE
        )
        AND "rfid"."santri_id" IS NULL
        OR sqlc.narg(is_santri) = TRUE
    )
    AND (
        (
            sqlc.narg(is_employee)::boolean IS NULL
            OR sqlc.narg(is_employee) = FALSE
        )
        AND "rfid"."employee_id" IS NULL
        OR sqlc.narg(is_employee) = TRUE
    )
ORDER BY
    "rfid"."id" ASC
LIMIT
    @limit_number OFFSET @offset_number;

-- name: UpdateRfid :one
UPDATE
    rfid
SET
    "uid" = @uid,
    "is_active" = @is_active,
    "santri_id" = sqlc.narg(santri_id),
    "employee_id" = sqlc.narg(employee_id)
WHERE
    "id" = @id RETURNING *;

-- name: GetRfidById :one
SELECT
    "rfid".*,
    "santri"."name" as "santri_name",
    "employee"."name" as "employee_name"
FROM
    rfid
    LEFT JOIN "santri" ON "rfid"."santri_id" = "santri"."id"
    LEFT JOIN "employee" ON "rfid"."employee_id" = "employee"."id"
WHERE
    "rfid"."id" = @id;

-- name: DeleteRfid :one
DELETE FROM
    rfid
WHERE
    "id" = @id RETURNING *;