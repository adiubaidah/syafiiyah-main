package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearSantriPermissionTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), "DELETE FROM santri_permission")
	require.NoError(t, err)
}

func createRandomSantriPermission(t *testing.T) (SantriPermission, Santri) {
	santri := createRandomSantri(t)
	startPermission := random.RandomTimeStamp().In(time.Local)
	endTime := startPermission.Add(time.Hour * 2)

	permissionTypes := []SantriPermissionType{SantriPermissionTypePermission, SantriPermissionTypeSick}

	arg := CreateSantriPermissionParams{
		SantriID:        santri.ID,
		StartPermission: pgtype.Timestamptz{Time: startPermission, Valid: true},
		Type:            permissionTypes[random.RandomInt(0, 1)],
		EndPermission:   pgtype.Timestamptz{Time: endTime, Valid: true},
		Excuse:          random.RandomString(100),
	}

	santriPermission, err := testStore.CreateSantriPermission(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, santriPermission)

	require.Equal(t, arg.SantriID, santriPermission.SantriID)
	require.Equal(t, arg.StartPermission.Time, santriPermission.StartPermission.Time)
	require.Equal(t, arg.EndPermission.Time, santriPermission.EndPermission.Time)
	require.Equal(t, arg.Excuse, santriPermission.Excuse)
	require.NotZero(t, santriPermission.ID)

	return santriPermission, santri
}

func TestCreateSantriPermission(t *testing.T) {
	clearSantriPermissionTable(t)
	clearSantriTable(t)
	createRandomSantriPermission(t)
}

func TestListSantriPermission(t *testing.T) {
	clearSantriPermissionTable(t)
	clearSantriTable(t)
	randomSantriPermission, santri := createRandomSantriPermission(t)
	for i := 0; i < 15; i++ {
		createRandomSantriPermission(t)
	}

	t.Run("list santri should contain santri name", func(t *testing.T) {
		arg := ListSantriPermissionsParams{
			Q:            pgtype.Text{String: santri.Name, Valid: true},
			OffsetNumber: 0,
			LimitNumber:  10,
		}
		santriPermissions, err := testStore.ListSantriPermissions(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, santriPermissions)

		require.Equal(t, santriPermissions[0].SantriID, santri.ID)
		require.Equal(t, santriPermissions[0].SantriName, santri.Name)
	})

	t.Run("list santri should contain santri permission type", func(t *testing.T) {
		arg := ListSantriPermissionsParams{
			Type:         NullSantriPermissionType{SantriPermissionType: SantriPermissionTypePermission, Valid: true},
			OffsetNumber: 0,
			LimitNumber:  10,
		}
		santriPermissions, err := testStore.ListSantriPermissions(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, santriPermissions)

		for _, santriPermission := range santriPermissions {
			require.Equal(t, SantriPermissionTypePermission, santriPermission.Type)
		}
	})

	t.Run("list santri should contain start_permission greater or equal than from_date", func(t *testing.T) {
		arg := ListSantriPermissionsParams{
			FromDate:     pgtype.Timestamptz{Time: randomSantriPermission.StartPermission.Time, Valid: true},
			OffsetNumber: 0,
			LimitNumber:  3,
		}
		santriPermissions, err := testStore.ListSantriPermissions(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, santriPermissions)

		for _, santriPermission := range santriPermissions {

			require.True(t, santriPermission.StartPermission.Time.After(arg.FromDate.Time) || santriPermission.StartPermission.Time.Equal(arg.FromDate.Time))
		}
	})
	t.Run("list santri should contain start_permission greater or equal than from_date", func(t *testing.T) {
		arg := ListSantriPermissionsParams{
			EndDate:      pgtype.Timestamptz{Time: randomSantriPermission.EndPermission.Time, Valid: true},
			OffsetNumber: 0,
			LimitNumber:  3,
		}
		santriPermissions, err := testStore.ListSantriPermissions(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, santriPermissions)

		for _, santriPermission := range santriPermissions {
			require.True(t, santriPermission.EndPermission.Time.Before(arg.EndDate.Time) || santriPermission.EndPermission.Time.Equal(arg.EndDate.Time))
		}
	})

}
