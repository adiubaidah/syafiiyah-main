package db

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearSantriPresenceTable(t *testing.T) {
	_, err := testQueries.db.Exec(context.Background(), `DELETE FROM "santri_presence"`)
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
		Type:         types[util.RandomInt(0, int64(n-1))],
		Notes:        pgtype.Text{Valid: false},
	}

	santriPresence, err := testQueries.CreateSantriPresence(context.Background(), arg)
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
	clearSantriPresenceTable(t)
	clearSantriTable(t)
	clearSantriSchedule(t)
	createRandomSantriPresence(t)
}

func TestListSantriPresence(t *testing.T) {
	clearSantriPresenceTable(t)
	santriPresence := createRandomSantriPresence(t)
	for i := 0; i < 30; i++ {
		createRandomSantriPresence(t)
	}

	t.Run("list santri should contain santri id", func(t *testing.T) {
		arg := ListSantriPresencesParams{
			SantriID:     pgtype.Int4{Int32: santriPresence.SantriID, Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
		}
		santriPresences, err := testQueries.ListSantriPresences(context.Background(), arg)
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
		santriPresences, err := testQueries.ListSantriPresences(context.Background(), arg)
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
		santriPresences, err := testQueries.ListSantriPresences(context.Background(), arg)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(santriPresences), 1)
	})

	t.Run("list santri pagination", func(t *testing.T) {
		testCases := []struct {
			name     string
			arg      ListSantriPresencesParams
			expected int
		}{
			{
				name: "Limit 5",
				arg: ListSantriPresencesParams{
					LimitNumber:  5,
					OffsetNumber: 0,
				},
				expected: 5,
			},
			{
				name: "Limit 5 Offset 5",
				arg: ListSantriPresencesParams{
					LimitNumber:  5,
					OffsetNumber: 5,
				},
				expected: 5,
			},
			{
				name: "Limit 5 Offset 10",
				arg: ListSantriPresencesParams{
					LimitNumber:  5,
					OffsetNumber: 10,
				},
				expected: 0,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				employees, err := testQueries.ListSantriPresences(context.Background(), tt.arg)
				require.NoError(t, err)
				require.Len(t, employees, tt.expected)
			})
		}
	})

}
