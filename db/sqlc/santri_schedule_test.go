package db

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearSantriSchedule(t *testing.T) {
	_, err := testQueries.db.Exec(context.Background(), "DELETE FROM santri_schedule")
	require.NoError(t, err)
}

func createRandomSantriSchedule(t *testing.T) SantriSchedule {
	startPresence, err := util.ConvertToPgxTime("07:00")
	require.NoError(t, err)
	startTime, err := util.ConvertToPgxTime("07:30")
	require.NoError(t, err)
	finishTime, err := util.ConvertToPgxTime("15:00")
	require.NoError(t, err)
	arg := CreateSantriScheduleParams{
		Name:          util.RandomString(10),
		Description:   pgtype.Text{Valid: false},
		StartPresence: startPresence,
		StartTime:     startTime,
		FinishTime:    finishTime,
	}
	santriSchedule, err := testQueries.CreateSantriSchedule(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, santriSchedule)

	require.Equal(t, arg.Name, santriSchedule.Name)
	require.Equal(t, arg.Description.String, santriSchedule.Description.String)
	require.Equal(t, arg.StartPresence, santriSchedule.StartPresence)
	require.Equal(t, arg.StartTime, santriSchedule.StartTime)
	require.Equal(t, arg.FinishTime, santriSchedule.FinishTime)
	return santriSchedule
}

func TestCreateSantriSchedule(t *testing.T) {
	createRandomSantriSchedule(t)
}

func TestListSantriSchedule(t *testing.T) {
	clearSantriSchedule(t)
	for i := 0; i < 10; i++ {
		createRandomSantriSchedule(t)
	}

	santriSchedules, err := testQueries.ListSantriSchedules(context.Background())
	require.NoError(t, err)
	require.Len(t, santriSchedules, 10)
}

func TestUpdateSantriSchedule(t *testing.T) {
	clearSantriSchedule(t)
	santriSchedule := createRandomSantriSchedule(t)

	startPresence, err := util.ConvertToPgxTime("06:45")
	require.NoError(t, err)
	startTime, err := util.ConvertToPgxTime("07:15")
	require.NoError(t, err)
	finishTime, err := util.ConvertToPgxTime("14:45")
	require.NoError(t, err)
	arg := UpdateSantriScheduleParams{
		ID:            santriSchedule.ID,
		Name:          util.RandomString(10),
		Description:   pgtype.Text{Valid: false},
		StartPresence: startPresence,
		StartTime:     startTime,
		FinishTime:    finishTime,
	}
	updatedSantriSchedule, err := testQueries.UpdateSantriSchedule(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSantriSchedule)

	require.Equal(t, arg.ID, updatedSantriSchedule.ID)
	require.Equal(t, arg.Name, updatedSantriSchedule.Name)
	require.Equal(t, arg.Description.String, updatedSantriSchedule.Description.String)
	require.Equal(t, arg.StartPresence, updatedSantriSchedule.StartPresence)
	require.Equal(t, arg.StartTime, updatedSantriSchedule.StartTime)
	require.Equal(t, arg.FinishTime, updatedSantriSchedule.FinishTime)
}

func TestDeleteSantriSchedule(t *testing.T) {
	clearSantriSchedule(t)
	santriSchedule := createRandomSantriSchedule(t)

	deletedSantriSchedule, err := testQueries.DeleteSantriSchedule(context.Background(), santriSchedule.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedSantriSchedule)

	require.Equal(t, santriSchedule.ID, deletedSantriSchedule.ID)
	require.Equal(t, santriSchedule.Name, deletedSantriSchedule.Name)
	require.Equal(t, santriSchedule.Description.String, deletedSantriSchedule.Description.String)
	require.Equal(t, santriSchedule.StartPresence, deletedSantriSchedule.StartPresence)
	require.Equal(t, santriSchedule.StartTime, deletedSantriSchedule.StartTime)
	require.Equal(t, santriSchedule.FinishTime, deletedSantriSchedule.FinishTime)
}
