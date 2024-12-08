package persistence

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

func (store *SQLStore) CreateHolidayWithDates(ctx context.Context, arg CreateHolidayParams, argsCreateDates []CreateHolidayDatesParams) (Holiday, error) {
	var createdHoliday Holiday

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		holiday, err := q.CreateHoliday(ctx, arg)
		if err != nil {
			return err
		}
		createdHoliday = holiday
		var args []CreateHolidayDatesParams
		for _, arg := range argsCreateDates {
			arg.HolidayID = holiday.ID
			args = append(args, arg)

		}
		_, err = q.CreateHolidayDates(ctx, args)

		return err
	})
	return createdHoliday, err
}
func (store *SQLStore) UpdateHolidayWithDates(ctx context.Context, holidayId int32, arg UpdateHolidayParams, argsCreateDates []CreateHolidayDatesParams) (Holiday, error) {
	var updateHoliday Holiday

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		holiday, err := q.UpdateHoliday(ctx, arg)
		if err != nil {
			return err
		}
		updateHoliday = holiday

		err = q.DeleteHolidayDateByHolidayId(ctx, holidayId)
		if err != nil {
			return err
		}

		var args []CreateHolidayDatesParams
		for _, arg := range argsCreateDates {
			arg.HolidayID = holiday.ID
			args = append(args, arg)

		}
		_, err = q.CreateHolidayDates(ctx, args)

		return err
	})
	return updateHoliday, err
}
