-- name: CreateSantriSchedule :one
INSERT INTO
    "santri_schedule" (
        "name",
        "description",
        "start_presence",
        "start_time",
        "finish_time"
    )
VALUES
    (
        @name,
        sqlc.narg(description),
        @start_presence::time,
        @start_time::time,
        @finish_time::time
    )
RETURNING *;

-- name: ListSantriSchedules :many
SELECT
    *
FROM
    "santri_schedule";

-- name: UpdateSantriSchedule :one
UPDATE
    "santri_schedule"
SET
    "name" = @name,
    "description" = sqlc.narg(description),
    "start_presence" = @start_presence::time,
    "start_time" = @start_time::time,
    "finish_time" = @finish_time::time
WHERE
    "id" = @id
RETURNING *;

-- name: DeleteSantriSchedule :one
DELETE FROM
    "santri_schedule"
WHERE
    "id" = @id
RETURNING *;