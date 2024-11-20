-- name: CreateSantriPermission :one
INSERT INTO
    "santri_permission" (
        santri_id,
        start_permission,
        end_permission,
        "type",
        excuse
    )
VALUES
    (
        @santri_id,
        @start_permission,
        sqlc.narg(end_permission),
        @type :: santri_permission_type,
        @excuse
    ) RETURNING *;

-- name: ListSantriPermissions :many
SELECT
    "santri_permission".*,
    "santri"."name" AS "santri_name"
FROM
    "santri_permission"
    INNER JOIN "santri" ON "santri_permission"."santri_id" = "santri"."id"
WHERE
    (sqlc.narg(q) :: text IS NULL
    OR "santri"."name" ILIKE '%' || sqlc.narg(q) || '%')
    AND (
        sqlc.narg(santri_id) :: integer IS NULL
        OR "santri_id" = sqlc.narg(santri_id) :: integer
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

-- name: GetSantriPermission :one
SELECT
    "santri_permission".*,
    "santri"."name" AS "santri_name"
FROM
    "santri_permission"
    INNER JOIN "santri" ON "santri_permission"."santri_id" = "santri"."id"
WHERE
    "santri_permission"."id" = @id;

-- name: UpdateSantriPermission :one
UPDATE
    "santri_permission"
SET
    "santri_id" = COALESCE(sqlc.narg(santri_id), santri_id),
    "start_permission" = COALESCE(sqlc.narg(start_permission), start_permission),
    "end_permission" = sqlc.narg(end_permission),
    "excuse" = COALESCE(sqlc.narg(excuse), excuse)
WHERE
    "id" = @id RETURNING *;

-- name: DeleteSantriPermission :one
DELETE FROM
    "santri_permission"
WHERE
    "id" = @id RETURNING *;