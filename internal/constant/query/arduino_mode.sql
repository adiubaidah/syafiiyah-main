
-- name: CreateArduinoModes :copyfrom
INSERT INTO
    "arduino_mode" (
        "mode",
        "input_topic",
        "acknowledgment_topic",
        "arduino_id"
    )
VALUES
    (
        @mode :: arduino_mode_type,
        @input_topic,
        @acknowledgement_topic,
        @arduino_id
    );

-- name: ListArduinoModes :many
SELECT
    *
FROM
    "arduino_mode"
WHERE
    arduino_id = @arduino_id;

-- name: UpdateArduinoMode :one
UPDATE
    "arduino_mode"
SET
    "mode" = COALESCE(sqlc.narg(mode)),
    "input_topic" = COALESCE(sqlc.narg(input_topic), input_topic),
    "acknowledgment_topic" = COALESCE(sqlc.narg(acknowledgement_topic), acknowledgement_topic)
WHERE
    "id" = @id RETURNING *;

-- name: DeleteArduinoModeByArduinoId :exec
DELETE FROM
    "arduino_mode"
WHERE
    "arduino_id" = @arduino_id;