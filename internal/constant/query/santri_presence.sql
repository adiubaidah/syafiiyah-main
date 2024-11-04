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
        @created_by :: presence_created_by,
        @santri_permission_id
    ) RETURNING *;

-- name: ListSantriPresences :many
SELECT
    *
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
        sqlc.narg(from_date) :: timestamp IS NULL
        OR "created_at" >= sqlc.narg(from_date) :: timestamp
    )
    AND (
        sqlc.narg(to_date) :: timestamp IS NULL
        OR "created_at" <= sqlc.narg(to_date) :: timestamp
    )
LIMIT
    @limit_number OFFSET @offset_number;

-- name: UpdateSantriPresence :one
UPDATE
    "santri_presence"
SET
    "schedule_id" = @schedule_id,
    "schedule_name" = @schedule_name,
    "type" = @type::presence_type,
    "santri_id" = @santri_id,
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