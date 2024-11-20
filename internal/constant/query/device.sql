-- name: CreateDevice :one
INSERT INTO
    "device" ("name")
VALUES
    (@name) RETURNING *;

-- name: ListDevices :many
SELECT
    "device"."id" AS "id",
    "device"."name" AS "name",
    "device_mode"."id" AS "device_mode.id",
    "device_mode"."mode" AS "device_mode.mode",
    "device_mode"."input_topic" AS "device_mode.input_topic",
    "device_mode"."acknowledgment_topic" AS "device_mode.acknowledgement_topic"
FROM
    "device"
LEFT JOIN
    "device_mode" ON "device"."id" = "device_mode"."device_id";

-- name: UpdateDevice :one
UPDATE
    "device"
SET
    "name" = COALESCE(sqlc.narg(name), name)
WHERE
    "id" = @id RETURNING *;

-- name: DeleteDevice :one
DELETE FROM
    "device"
WHERE
    "id" = @id RETURNING *;