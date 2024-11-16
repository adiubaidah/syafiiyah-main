-- name: CreateArduino :one
INSERT INTO
    "arduino" ("name")
VALUES
    (@name) RETURNING *;

-- name: ListArduinos :many
SELECT
    "arduino"."id" AS "id",
    "arduino"."name" AS "name",
    "arduino_mode"."id" AS "arduino_mode.id",
    "arduino_mode"."mode" AS "arduino_mode.mode",
    "arduino_mode"."input_topic" AS "arduino_mode.input_topic",
    "arduino_mode"."acknowledgment_topic" AS "arduino_mode.acknowledgement_topic"
FROM
    "arduino"
LEFT JOIN
    "arduino_mode" ON "arduino"."id" = "arduino_mode"."arduino_id";

-- name: UpdateArduino :one
UPDATE
    "arduino"
SET
    "name" = COALESCE(sqlc.narg(name), name)
WHERE
    "id" = @id RETURNING *;

-- name: DeleteArduino :one
DELETE FROM
    "arduino"
WHERE
    "id" = @id RETURNING *;