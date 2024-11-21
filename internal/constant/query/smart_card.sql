-- name: CreateSmartCard :one
INSERT INTO
    smart_card ("uid", "is_active", "santri_id", "employee_id")
VALUES
    (
        @uid,
        @is_active,
        sqlc.narg(santri_id),
        sqlc.narg(employee_id)
    ) RETURNING *;

-- name: ListSmartCards :many
SELECT
    "smart_card".*,
    "santri"."name" as "santri_name",
    "employee"."name" as "employee_name"
FROM
    smart_card
    LEFT JOIN "santri" ON "smart_card"."santri_id" = "santri"."id"
    LEFT JOIN "employee" ON "smart_card"."employee_id" = "employee"."id"
WHERE
    (
        sqlc.narg(q) :: text IS NULL
        OR "uid" ILIKE '%' || sqlc.narg(q) || '%'
        OR "santri"."name" ILIKE '%' || sqlc.narg(q) || '%'
        OR "employee"."name" ILIKE '%' || sqlc.narg(q) || '%'
    )
    AND (
        sqlc.narg(is_active)::boolean IS NULL
        OR "smart_card"."is_active" = sqlc.narg(is_active)
    )
    AND (
        (
            sqlc.narg(is_santri)::boolean IS NULL
            OR sqlc.narg(is_santri) = FALSE
        )
        AND "smart_card"."santri_id" IS NULL
        OR sqlc.narg(is_santri) = TRUE
    )
    AND (
        (
            sqlc.narg(is_employee)::boolean IS NULL
            OR sqlc.narg(is_employee) = FALSE
        )
        AND "smart_card"."employee_id" IS NULL
        OR sqlc.narg(is_employee) = TRUE
    )
ORDER BY
    "smart_card"."id" ASC
LIMIT
    @limit_number OFFSET @offset_number;

-- name: CountSmartCards :one
SELECT
    COUNT(*) as "count"
FROM
    smart_card
WHERE
    (
        sqlc.narg(q) :: text IS NULL
        OR "uid" ILIKE '%' || sqlc.narg(q) || '%'
        OR "santri"."name" ILIKE '%' || sqlc.narg(q) || '%'
        OR "employee"."name" ILIKE '%' || sqlc.narg(q) || '%'
    )
    AND (
        sqlc.narg(is_active)::boolean IS NULL
        OR "smart_card"."is_active" = sqlc.narg(is_active)
    )
    AND (
        (
            sqlc.narg(is_santri)::boolean IS NULL
            OR sqlc.narg(is_santri) = FALSE
        )
        AND "smart_card"."santri_id" IS NULL
        OR sqlc.narg(is_santri) = TRUE
    )
    AND (
        (
            sqlc.narg(is_employee)::boolean IS NULL
            OR sqlc.narg(is_employee) = FALSE
        )
        AND "smart_card"."employee_id" IS NULL
        OR sqlc.narg(is_employee) = TRUE
    );

-- name: UpdateSmartCard :one
UPDATE
    smart_card
SET
    "uid" = COALESCE(sqlc.narg(uid), uid),
    "is_active" = COALESCE(sqlc.narg(is_active), is_active),
    "santri_id" = sqlc.narg(santri_id),
    "employee_id" = sqlc.narg(employee_id)
WHERE
    "id" = @id RETURNING *;

-- name: GetSmartCard :one
SELECT
    "smart_card".*,
    "santri"."name" as "santri_name",
    "employee"."name" as "employee_name"
FROM
    smart_card
    LEFT JOIN "santri" ON "smart_card"."santri_id" = "santri"."id"
    LEFT JOIN "employee" ON "smart_card"."employee_id" = "employee"."id"
WHERE
    "smart_card"."id" = @id;

-- name: DeleteSmartCard :one
DELETE FROM
    smart_card
WHERE
    "id" = @id RETURNING *;