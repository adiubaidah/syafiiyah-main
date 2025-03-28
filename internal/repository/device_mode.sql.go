// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: device_mode.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateDeviceModesParams struct {
	Mode                 DeviceModeType `db:"mode"`
	InputTopic           string         `db:"input_topic"`
	AcknowledgementTopic string         `db:"acknowledgement_topic"`
	DeviceID             int32          `db:"device_id"`
}

const deleteDeviceModeByDeviceId = `-- name: DeleteDeviceModeByDeviceId :exec
DELETE FROM
    "device_mode"
WHERE
    "device_id" = $1
`

func (q *Queries) DeleteDeviceModeByDeviceId(ctx context.Context, deviceID int32) error {
	_, err := q.db.Exec(ctx, deleteDeviceModeByDeviceId, deviceID)
	return err
}

const listDeviceModes = `-- name: ListDeviceModes :many
SELECT
    id, mode, input_topic, acknowledgment_topic, device_id
FROM
    "device_mode"
WHERE
    device_id = $1
`

func (q *Queries) ListDeviceModes(ctx context.Context, deviceID int32) ([]DeviceMode, error) {
	rows, err := q.db.Query(ctx, listDeviceModes, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []DeviceMode{}
	for rows.Next() {
		var i DeviceMode
		if err := rows.Scan(
			&i.ID,
			&i.Mode,
			&i.InputTopic,
			&i.AcknowledgmentTopic,
			&i.DeviceID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateDeviceMode = `-- name: UpdateDeviceMode :one
UPDATE
    "device_mode"
SET
    "mode" = COALESCE($1),
    "input_topic" = COALESCE($2, input_topic),
    "acknowledgment_topic" = COALESCE($3, acknowledgement_topic)
WHERE
    "id" = $4 RETURNING id, mode, input_topic, acknowledgment_topic, device_id
`

type UpdateDeviceModeParams struct {
	Mode                 NullDeviceModeType `db:"mode"`
	InputTopic           pgtype.Text        `db:"input_topic"`
	AcknowledgementTopic pgtype.Text        `db:"acknowledgement_topic"`
	ID                   int32              `db:"id"`
}

func (q *Queries) UpdateDeviceMode(ctx context.Context, arg UpdateDeviceModeParams) (DeviceMode, error) {
	row := q.db.QueryRow(ctx, updateDeviceMode,
		arg.Mode,
		arg.InputTopic,
		arg.AcknowledgementTopic,
		arg.ID,
	)
	var i DeviceMode
	err := row.Scan(
		&i.ID,
		&i.Mode,
		&i.InputTopic,
		&i.AcknowledgmentTopic,
		&i.DeviceID,
	)
	return i, err
}
