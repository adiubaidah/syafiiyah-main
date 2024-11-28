-- name: CreateSantriPresence :one
INSERT INTO
    "santri_presence" (
        "schedule_id",
        "schedule_name",
        "type",
        "santri_id",
        "notes",
        "created_by",
        "santri_permission_id"
    )
VALUES
    (
        @schedule_id,
        @schedule_name,
        @type :: presence_type,
        @santri_id,
        @notes,
        @created_by :: presence_created_by_type,
        @santri_permission_id
    ) RETURNING *;

-- name: CreateSantriPresences :copyfrom
INSERT INTO
    "santri_presence" (
        "schedule_id",
        "schedule_name",
        "type",
        "santri_id",
        "notes",
        "created_at",
        "created_by",
        "santri_permission_id"
    )
VALUES
    (
        @schedule_id,
        @schedule_name,
        @type::presence_type,
        @santri_id,
        @notes,
        @created_at,
        @created_by :: presence_created_by_type,
        @santri_permission_id
    );

-- name: ListSantriPresences :many
SELECT
    "santri_presence".*,
    "santri"."name" AS "santri_name"
FROM
    "santri_presence"
    INNER JOIN "santri" ON "santri_presence"."santri_id" = "santri"."id"
WHERE
    (
        sqlc.narg(santri_id) :: integer IS NULL
        OR "santri_id" = sqlc.narg(santri_id) :: integer
    )
    AND (
        sqlc.narg(q) :: text IS NULL
        OR "santri"."name" ILIKE '%' || sqlc.narg(q) || '%'
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
    "santri_presence"."id" DESC
LIMIT
    @limit_number OFFSET @offset_number;


-- name: CountSantriPresences :one
SELECT
    COUNT(*)
FROM
    "santri_presence"
    INNER JOIN "santri" ON "santri_presence"."santri_id" = "santri"."id"
WHERE
    (
        sqlc.narg(santri_id) :: integer IS NULL
        OR "santri_id" = sqlc.narg(santri_id) :: integer
    )
    AND (
        sqlc.narg(q) :: text IS NULL
        OR "santri"."name" ILIKE '%' || sqlc.narg(q) || '%'
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

-- name: ListAbsentSantri :many
SELECT 
    "santri"."id", "santri"."name"
FROM
    "santri"
WHERE
    NOT EXISTS (
        SELECT
            1
        FROM
            "santri_presence"
        WHERE
            "santri_presence"."santri_id" = "santri"."id"
            AND DATE("santri_presence"."created_at") = sqlc.narg(date)::date
            AND "santri_presence"."schedule_id" = sqlc.narg(schedule_id)::integer
    );

-- name: UpdateSantriPresence :one
UPDATE
    "santri_presence"
SET
    "schedule_id" = COALESCE(sqlc.narg(schedule_id), schedule_id),
    "schedule_name" = COALESCE(sqlc.narg(schedule_name), schedule_name),
    "type" = COALESCE(sqlc.narg(type)::presence_type, type),
    "santri_id" = COALESCE(sqlc.narg(santri_id), santri_id),
    "notes" = sqlc.narg(notes),
    "santri_permission_id" = @santri_permission_id
WHERE
    "id" = @id
RETURNING *;

-- name: DeleteSantriPresence :one
DELETE FROM
    "santri_presence"
WHERE
    "id" = @id
RETURNING *;