package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearSantriPresenceTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), `DELETE FROM "santri_presence"`)
	require.NoError(t, err)
}

func createRandomSantriPresence(t *testing.T) SantriPresence {
	santri := createRandomSantri(t)
	schedule := createRandomSantriSchedule(t)
	types := []PresenceType{PresenceTypeAlpha, PresenceTypeLate, PresenceTypePermission, PresenceTypePresent, PresenceTypeSick}
	n := len(types)
	arg := CreateSantriPresenceParams{
		SantriID:     santri.ID,
		ScheduleID:   schedule.ID,
		ScheduleName: schedule.Name,
		Type:         types[random.RandomInt(0, int64(n-1))],
		Notes:        pgtype.Text{Valid: false},
		CreatedBy:    PresenceCreatedByTypeTap,
	}

	santriPresence, err := testStore.CreateSantriPresence(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, santriPresence)

	require.Equal(t, arg.SantriID, santriPresence.SantriID)
	require.Equal(t, arg.ScheduleID, santriPresence.ScheduleID)
	require.Equal(t, arg.ScheduleName, santriPresence.ScheduleName)
	require.Equal(t, arg.Type, santriPresence.Type)
	require.Equal(t, arg.Notes.String, santriPresence.Notes.String)
	require.NotZero(t, santriPresence.ID)

	return santriPresence
}

func TestCreateSantriPresence(t *testing.T) {
	clearSantriPermissionTable(t)
	clearSantriPresenceTable(t)
	clearSantriTable(t)
	createRandomSantriPresence(t)
}

func TestListSantriPresence(t *testing.T) {
	clearSantriPresenceTable(t)
	santriPresence := createRandomSantriPresence(t)
	for i := 0; i < 10; i++ {
		createRandomSantriPresence(t)
	}

	t.Run("list santri should contain santri id", func(t *testing.T) {
		arg := ListSantriPresencesParams{
			SantriID:     pgtype.Int4{Int32: santriPresence.SantriID, Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
		}
		santriPresences, err := testStore.ListSantriPresences(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, santriPresences, 1)
		require.Equal(t, santriPresences[0].ID, santriPresence.ID)
	})

	t.Run("list santri should contain schedule id", func(t *testing.T) {
		arg := ListSantriPresencesParams{
			ScheduleID:   pgtype.Int4{Int32: santriPresence.ScheduleID, Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
		}
		santriPresences, err := testStore.ListSantriPresences(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, santriPresences, 1)
		require.Equal(t, santriPresences[0].ID, santriPresence.ID)
	})

	t.Run("list santri should match of type presence santri", func(t *testing.T) {
		arg := ListSantriPresencesParams{
			Type:         NullPresenceType{PresenceType: santriPresence.Type, Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
		}
		santriPresences, err := testStore.ListSantriPresences(context.Background(), arg)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(santriPresences), 1)
	})

	t.Run("list santri pagination", func(t *testing.T) {
		testCases := []struct {
			name        string
			arg         ListSantriPresencesParams
			lenExpected int
		}{
			{
				name: "Limit 5",
				arg: ListSantriPresencesParams{
					LimitNumber:  5,
					OffsetNumber: 0,
				},
				lenExpected: 5,
			},
			{
				name: "Limit 5 Offset 5",
				arg: ListSantriPresencesParams{
					LimitNumber:  5,
					OffsetNumber: 5,
				},
				lenExpected: 5,
			},
			{
				name: "Limit 5 Offset 10",
				arg: ListSantriPresencesParams{
					LimitNumber:  5,
					OffsetNumber: 10,
				},
				lenExpected: 1,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				employees, err := testStore.ListSantriPresences(context.Background(), tt.arg)
				require.NoError(t, err)
				require.Len(t, employees, tt.lenExpected)
			})
		}
	})

}

func TestListAbsentSantriPresence(t *testing.T) {
	clearSantriScheduleTable(t)
	clearSantriPresenceTable(t)
	clearSantriTable(t)

	santri := createRandomSantri(t)
	schedule := createRandomSantriSchedule(t)
	scheduleTime := time.Unix(0, schedule.StartPresence.Microseconds*int64(time.Microsecond))

	listAbsent, err := testStore.ListAbsentSantri(context.Background(), ListAbsentSantriParams{
		Date: pgtype.Date{Time: scheduleTime, Valid: true},
	})

	require.NoError(t, err)
	require.NotNil(t, listAbsent)

	require.Equal(t, listAbsent[0].ID, santri.ID)
	require.Equal(t, listAbsent[0].Name, santri.Name)
}

func TestUpdateSantriPresence(t *testing.T) {
	clearSantriPresenceTable(t)
	clearSantriTable(t)
	santriPresence := createRandomSantriPresence(t)
	arg := UpdateSantriPresenceParams{
		ID:                 santriPresence.ID,
		ScheduleID:         pgtype.Int4{Valid: false},
		ScheduleName:       pgtype.Text{Valid: false},
		Type:               NullPresenceType{PresenceType: PresenceTypeAlpha, Valid: true},
		SantriPermissionID: pgtype.Int4{Valid: false},
		SantriID:           pgtype.Int4{Int32: santriPresence.SantriID, Valid: true},
		Notes:              pgtype.Text{Valid: true, String: "Example notes"},
	}

	updatedSantriPresence, err := testStore.UpdateSantriPresence(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSantriPresence)

	require.NotEqual(t, arg.ScheduleID.Int32, updatedSantriPresence.ScheduleID)
	require.NotEqual(t, arg.ScheduleName.String, updatedSantriPresence.ScheduleName)
	require.Equal(t, arg.Type.PresenceType, updatedSantriPresence.Type)
	require.Equal(t, arg.SantriPermissionID.Int32, updatedSantriPresence.SantriPermissionID.Int32)
	require.Equal(t, arg.SantriID.Int32, updatedSantriPresence.SantriID)
	require.Equal(t, santriPresence.ID, updatedSantriPresence.ID)
	require.Equal(t, arg.Notes.String, updatedSantriPresence.Notes.String)
}

func TestDeleteSantriPresence(t *testing.T) {
	clearSantriPresenceTable(t)
	clearSantriTable(t)
	santriPresence := createRandomSantriPresence(t)

	deletedSantriPresence, err := testStore.DeleteSantriPresence(context.Background(), santriPresence.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedSantriPresence)

	require.Equal(t, santriPresence.ID, deletedSantriPresence.ID)
}
