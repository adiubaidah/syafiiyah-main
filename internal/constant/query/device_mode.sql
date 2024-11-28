
-- name: CreateDeviceModes :copyfrom
INSERT INTO
    "device_mode" (
        "mode",
        "input_topic",
        "acknowledgment_topic",
        "device_id"
    )
VALUES
    (
        @mode :: device_mode_type,
        @input_topic,
        @acknowledgement_topic,
        @device_id
    );

-- name: ListDeviceModes :many
SELECT
    *
FROM
    "device_mode"
WHERE
    device_id = @device_id;

-- name: UpdateDeviceMode :one
UPDATE
    "device_mode"
SET
    "mode" = COALESCE(sqlc.narg(mode)),
    "input_topic" = COALESCE(sqlc.narg(input_topic), input_topic),
    "acknowledgment_topic" = COALESCE(sqlc.narg(acknowledgement_topic), acknowledgement_topic)
WHERE
    "id" = @id RETURNING *;

-- name: DeleteDeviceModeByDeviceId :exec
DELETE FROM
    "device_mode"
WHERE
    "device_id" = @device_id;