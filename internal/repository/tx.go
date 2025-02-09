package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (store *SQLStore) CreateDeviceWithModes(ctx context.Context, arduinoName string, modeParams []CreateDeviceModesParams) (Device, error) {
	var createdDevice Device

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		device, err := q.CreateDevice(ctx, arduinoName)
		if err != nil {
			return err
		}
		createdDevice = device

		for i := range modeParams {
			modeParams[i].DeviceID = device.ID
		}
		_, err = q.CreateDeviceModes(ctx, modeParams)
		return err

	})
	return createdDevice, err
}

func (store *SQLStore) UpdateDeviceWithModes(ctx context.Context, deviceID int32,
	name string,
	modeParams []CreateDeviceModesParams) (Device, error) {
	var updatedArduino Device

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		device, err := q.UpdateDevice(ctx, UpdateDeviceParams{
			ID:   deviceID,
			Name: pgtype.Text{String: name, Valid: name != ""},
		})
		if err != nil {
			return err
		}

		updatedArduino = device

		err = q.DeleteDeviceModeByDeviceId(ctx, deviceID)
		if err != nil {
			return err
		}

		for i := range modeParams {
			modeParams[i].DeviceID = device.ID
		}

		_, err = q.CreateDeviceModes(ctx, modeParams)
		return err

	})
	return updatedArduino, err
}
