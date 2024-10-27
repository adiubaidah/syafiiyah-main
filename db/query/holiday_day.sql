-- name: CreateHolidayDay :one
INSERT INTO
    "holiday_day" ("date", "holiday_id")
VALUES
    (@date, @holiday_id) RETURNING *;

-- name: ListHolidayDays :many
SELECT
    *
FROM
    "holiday_day"
    INNER JOIN "holiday" ON "holiday_day"."holiday_id" = "holiday"."id"
WHERE
    "date" BETWEEN @from_date AND @to_date
    AND (
        sqlc.narg(holiday_id)::integer IS NULL
        OR "holiday_id" = @holiday_id
    )
    AND (
        sqlc.narg(holiday_name)::text IS NULL
        OR "holiday"."name" ILIKE '%' || sqlc.narg(holiday_name) || '%'
    )
    LIMIT @limit_number OFFSET @offset_number;

-- name: DeleteHolidayDay :one
DELETE FROM
    "holiday_day"
WHERE
    "id" = @id RETURNING *;