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
        @start_presence :: time,
        @start_time :: time,
        @finish_time :: time
    ) RETURNING *;

-- name: ExistBetweenSantriSchedule :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            "santri_schedule"
        WHERE
            "start_time" <= @start_time :: time
            AND "finish_time" >= @finish_time :: time
    ) AS "exist";

-- name: ListSantriSchedules :many
SELECT
    *
FROM
    "santri_schedule"
ORDER BY
    "start_time" ASC;

-- name: GetLastSantriSchedule :one
SELECT
    *
FROM
    "santri_schedule"
WHERE
    start_time = (
        SELECT
            MAX(start_time)
        FROM
            "santri_schedule"
    );

-- name: UpdateSantriSchedule :one
UPDATE
    "santri_schedule"
SET
    "name" = @name,
    "description" = sqlc.narg(description),
    "start_presence" = COALESCE(sqlc.narg(start_presence) :: time, start_presence),
    "start_time" = COALESCE(sqlc.narg(start_time) :: time, start_time),
    "finish_time" = COALESCE(sqlc.narg(finish_time) :: time, finish_time)
WHERE
    "id" = @id RETURNING *;

-- name: DeleteSantriSchedule :one
DELETE FROM
    "santri_schedule"
WHERE
    "id" = @id RETURNING *;