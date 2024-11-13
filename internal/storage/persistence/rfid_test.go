package persistence

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearRfidTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), `DELETE FROM "rfid"`)
	require.NoError(t, err)
}

func createRandomRfidWithSantri(t *testing.T) (Rfid, Santri) {
	santri := createRandomSantri(t)
	arg := CreateRfidParams{
		Uid:      random.RandomString(12),
		IsActive: random.RandomBool(),
		SantriID: pgtype.Int4{Int32: santri.ID, Valid: true},
	}
	rfid, err := testStore.CreateRfid(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, rfid)

	require.Equal(t, arg.Uid, rfid.Uid)
	require.Equal(t, arg.IsActive, rfid.IsActive)
	require.Equal(t, arg.SantriID, rfid.SantriID)
	require.Equal(t, arg.EmployeeID, rfid.EmployeeID)

	require.NotZero(t, rfid.ID)
	require.NotZero(t, rfid.CreatedAt)

	return rfid, santri
}

func createRandomRfidWithEmployee(t *testing.T) (Rfid, Employee) {
	employee := createRandomEmployee(t)
	arg := CreateRfidParams{
		Uid:        random.RandomString(12),
		IsActive:   random.RandomBool(),
		EmployeeID: pgtype.Int4{Int32: employee.ID, Valid: true},
	}
	rfid, err := testStore.CreateRfid(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, rfid)

	require.Equal(t, arg.Uid, rfid.Uid)
	require.Equal(t, arg.IsActive, rfid.IsActive)
	require.Equal(t, arg.SantriID, rfid.SantriID)
	require.Equal(t, arg.EmployeeID, rfid.EmployeeID)

	require.NotZero(t, rfid.ID)
	require.NotZero(t, rfid.CreatedAt)

	return rfid, employee
}

func TestCreateRfid(t *testing.T) {
	clearRfidTable(t)
	clearSantriTable(t)
	createRandomRfidWithSantri(t)
}

func TestListRfids(t *testing.T) {
	clearRfidTable(t)
	clearSantriTable(t)
	clearEmployeeTable(t)

	randomSantriRfid, santri := createRandomRfidWithSantri(t)
	randomEmployeeRfid, employee := createRandomRfidWithEmployee(t)

	for i := 0; i < 10; i++ {
		createRandomRfidWithSantri(t)
		createRandomRfidWithEmployee(t)
	}

	t.Run("list all rfids should match santri", func(t *testing.T) {
		arg := ListRfidParams{
			Q:            pgtype.Text{String: santri.Name[:3], Valid: true},
			OffsetNumber: 0,
			LimitNumber:  10,
			IsSantri:     pgtype.Bool{Bool: true, Valid: true},
		}
		rfids, err := testStore.ListRfid(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, rfids)

		for _, rfid := range rfids {
			require.NotEmpty(t, rfid)
			require.Equal(t, randomSantriRfid.Uid, rfid.Uid)
			require.Equal(t, randomSantriRfid.IsActive, rfid.IsActive)
			require.Equal(t, randomSantriRfid.SantriID, rfid.SantriID)

			//rfid.employee_id should be null
			require.False(t, rfid.EmployeeID.Valid)
			require.True(t, rfid.SantriID.Valid)
			require.Equal(t, santri.Name, rfid.SantriName.String)
		}
	})

	t.Run("list all rfids should match employee", func(t *testing.T) {
		arg := ListRfidParams{
			Q:            pgtype.Text{String: employee.Name[:3], Valid: true},
			OffsetNumber: 0,
			LimitNumber:  10,
			IsEmployee:   pgtype.Bool{Bool: true, Valid: true},
		}
		rfids, err := testStore.ListRfid(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, rfids)

		for _, rfid := range rfids {
			require.NotEmpty(t, rfid)
			require.Equal(t, randomEmployeeRfid.Uid, rfid.Uid)
			require.Equal(t, randomEmployeeRfid.IsActive, rfid.IsActive)
			require.Equal(t, randomEmployeeRfid.EmployeeID, rfid.EmployeeID)

			//rfid.santri_id should be null
			require.False(t, rfid.SantriID.Valid)
			require.True(t, rfid.EmployeeID.Valid)
			require.Equal(t, employee.Name, rfid.EmployeeName.String)
		}
	})

}

func TestUpdateRfid(t *testing.T) {
	clearRfidTable(t)
	clearSantriTable(t)
	clearEmployeeTable(t)

	rfid, _ := createRandomRfidWithSantri(t)

	arg := UpdateRfidParams{
		ID:         rfid.ID,
		Uid:        pgtype.Text{String: random.RandomString(12), Valid: true},
		IsActive:   pgtype.Bool{Bool: random.RandomBool(), Valid: true},
		SantriID:   pgtype.Int4{Int32: 0, Valid: false},
		EmployeeID: pgtype.Int4{Int32: 0, Valid: false},
	}
	updatedRfid, err := testStore.UpdateRfid(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedRfid)

	require.Equal(t, arg.ID, updatedRfid.ID)
	require.Equal(t, arg.Uid.String, updatedRfid.Uid)
	require.Equal(t, arg.IsActive.Bool, updatedRfid.IsActive)
	require.Equal(t, arg.SantriID.Int32, updatedRfid.SantriID.Int32)
	require.Equal(t, arg.EmployeeID.Int32, updatedRfid.EmployeeID.Int32)
}

func TestDeleteRfid(t *testing.T) {
	clearRfidTable(t)
	clearSantriTable(t)
	clearEmployeeTable(t)

	rfid, _ := createRandomRfidWithSantri(t)

	deletedRfid, err := testStore.DeleteRfid(context.Background(), rfid.ID)
	require.NoError(t, err)

	require.NotEmpty(t, deletedRfid)
}
