-- name: CreateHoliday :one
INSERT INTO
    "holiday" ("name", "color", "description")
VALUES
    (@name, @color, @description) RETURNING *;

-- name: ListHolidays :many
SELECT
    "holiday".*,
    "holiday_date"."id" AS "holiday_date_id",
    "holiday_date"."date" AS "holiday_date"
FROM
    "holiday"
    LEFT JOIN "holiday_date" ON "holiday"."id" = "holiday_date"."holiday_id"
WHERE
    (
        sqlc.narg(q)::text IS NULL
        OR "holiday"."name" ILIKE sqlc.arg(q)
    )
    AND
    (
        sqlc.narg(month)::integer IS NULL
        OR EXTRACT(
            MONTH
            FROM
                "holiday_date"."date"
        ) = CAST(sqlc.arg(month) AS INTEGER)
    )
    AND (
        sqlc.narg(year)::integer IS NULL
        OR EXTRACT(
            YEAR
            FROM
                "holiday_date"."date"
        ) = COALESCE(sqlc.arg(year), EXTRACT(YEAR FROM CURRENT_DATE))
    )
ORDER BY
    "holiday_date"."date" ASC;

-- name: UpdateHoliday :one
UPDATE
    "holiday"
SET
    "name" = COALESCE(sqlc.narg(name), "name"),
    "color" = sqlc.narg(color),
    "description" = sqlc.narg(description)
WHERE
    "id" = @id RETURNING *;

-- name: DeleteHoliday :one
DELETE FROM
    "holiday"
WHERE
    "id" = @id RETURNING *;