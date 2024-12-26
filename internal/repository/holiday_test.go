package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearHoliday(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), `DELETE FROM "holiday"`)
	require.NoError(t, err)
}

func createRandomHoliday(t *testing.T) Holiday {
	arg := CreateHolidayParams{
		Name:        random.RandomString(8),
		Color:       pgtype.Text{String: random.RandomString(7), Valid: true},
		Description: pgtype.Text{String: random.RandomString(50), Valid: true},
	}

	holiday, err := testStore.CreateHoliday(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, holiday)

	require.Equal(t, arg.Name, holiday.Name)
	require.Equal(t, arg.Description.String, holiday.Description.String)
	return holiday
}

func createHolidayDateWithDate(t *testing.T, holidayID int32, date string) {
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, date)
	require.NoError(t, err)

	var args []CreateHolidayDatesParams

	args = append(args, CreateHolidayDatesParams{
		HolidayID: holidayID,
		Date:      pgtype.Date{Time: parsedDate, Valid: true},
	})

	_, err = testStore.CreateHolidayDates(context.Background(), args)
	require.NoError(t, err)
}

func createHolidayDate(t *testing.T, holidayID int32) {
	var args []CreateHolidayDatesParams

	for i := 0; i < 5; i++ {
		args = append(args, CreateHolidayDatesParams{
			HolidayID: holidayID,
			Date:      pgtype.Date{Time: random.RandomTimeStamp(), Valid: true},
		})
	}

	affectedRows, err := testStore.CreateHolidayDates(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, int64(len(args)), affectedRows)

}

func TestCreateHoliday(t *testing.T) {
	clearHoliday(t)
	randomHoliday := createRandomHoliday(t)
	createHolidayDate(t, randomHoliday.ID)
}

func TestListHoliday(t *testing.T) {
	clearHoliday(t)
	// Create holiday for December 2023
	holidayDec := createRandomHoliday(t)
	createHolidayDateWithDate(t, holidayDec.ID, "2023-12-25")

	// Create holiday for January 2024
	holidayJan := createRandomHoliday(t)
	createHolidayDateWithDate(t, holidayJan.ID, "2024-01-01")

	t.Run("test for december 2023", func(t *testing.T) {
		holidays, err := testStore.ListHolidays(context.Background(), ListHolidaysParams{
			Month: pgtype.Int4{Int32: 12, Valid: true},
			Year:  pgtype.Int4{Int32: 2023, Valid: true},
		})
		require.NoError(t, err)
		require.Len(t, holidays, 1)
		require.Equal(t, holidayDec.ID, holidays[0].ID)
	})

	t.Run("test for january 2024", func(t *testing.T) {
		holidays, err := testStore.ListHolidays(context.Background(), ListHolidaysParams{
			Month: pgtype.Int4{Int32: 1, Valid: true},
			Year:  pgtype.Int4{Int32: 2024, Valid: true},
		})
		require.NoError(t, err)
		require.Len(t, holidays, 1)
		require.Equal(t, holidayJan.ID, holidays[0].ID)
	})

	t.Run("test for no holidays", func(t *testing.T) {
		holidays, err := testStore.ListHolidays(context.Background(), ListHolidaysParams{
			Month: pgtype.Int4{Int32: 5, Valid: true},
			Year:  pgtype.Int4{Int32: 2024, Valid: true},
		})
		require.NoError(t, err)
		require.Len(t, holidays, 0)
	})
}

func TestDeleteHoliday(t *testing.T) {
	clearHoliday(t)
	holiday := createRandomHoliday(t)
	deletedHoliday, err := testStore.DeleteHoliday(context.Background(), holiday.ID)

	require.NoError(t, err)
	require.NotEmpty(t, deletedHoliday)
}
