package persistence

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearSmartCardTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), `DELETE FROM "smart_card"`)
	require.NoError(t, err)
}

func createRandomSmartCardWithSantri(t *testing.T) (SmartCard, Santri) {
	santri := createRandomSantri(t)
	arg := CreateSmartCardParams{
		Uid:      random.RandomString(12),
		IsActive: random.RandomBool(),
		SantriID: pgtype.Int4{Int32: santri.ID, Valid: true},
	}
	smartCard, err := testStore.CreateSmartCard(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, smartCard)

	require.Equal(t, arg.Uid, smartCard.Uid)
	require.Equal(t, arg.IsActive, smartCard.IsActive)
	require.Equal(t, arg.SantriID, smartCard.SantriID)
	require.Equal(t, arg.EmployeeID, smartCard.EmployeeID)

	require.NotZero(t, smartCard.ID)
	require.NotZero(t, smartCard.CreatedAt)

	return smartCard, santri
}

func createRandomRfidWithEmployee(t *testing.T) (SmartCard, Employee) {
	employee := createRandomEmployee(t)
	arg := CreateSmartCardParams{
		Uid:        random.RandomString(12),
		IsActive:   random.RandomBool(),
		EmployeeID: pgtype.Int4{Int32: employee.ID, Valid: true},
	}
	smartCard, err := testStore.CreateSmartCard(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, smartCard)

	require.Equal(t, arg.Uid, smartCard.Uid)
	require.Equal(t, arg.IsActive, smartCard.IsActive)
	require.Equal(t, arg.SantriID, smartCard.SantriID)
	require.Equal(t, arg.EmployeeID, smartCard.EmployeeID)

	require.NotZero(t, smartCard.ID)
	require.NotZero(t, smartCard.CreatedAt)

	return smartCard, employee
}

func TestCreateSmartCard(t *testing.T) {
	clearSmartCardTable(t)
	clearSantriTable(t)
	createRandomSmartCardWithSantri(t)
}

func TestListSmartCards(t *testing.T) {
	clearSmartCardTable(t)
	clearSantriTable(t)
	clearEmployeeTable(t)

	randomSantriRfid, santri := createRandomSmartCardWithSantri(t)
	randomEmployeeRfid, employee := createRandomRfidWithEmployee(t)

	for i := 0; i < 10; i++ {
		createRandomSmartCardWithSantri(t)
		createRandomRfidWithEmployee(t)
	}

	t.Run("list all rfids should match santri", func(t *testing.T) {
		arg := ListSmartCardsParams{
			Q:            pgtype.Text{String: santri.Name[:3], Valid: true},
			OffsetNumber: 0,
			LimitNumber:  10,
			IsSantri:     pgtype.Bool{Bool: true, Valid: true},
		}
		rfids, err := testStore.ListSmartCards(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, rfids)

		for _, smartCard := range rfids {
			require.NotEmpty(t, smartCard)
			require.Equal(t, randomSantriRfid.Uid, smartCard.Uid)
			require.Equal(t, randomSantriRfid.IsActive, smartCard.IsActive)
			require.Equal(t, randomSantriRfid.SantriID, smartCard.SantriID)

			//smartCard.employee_id should be null
			require.False(t, smartCard.EmployeeID.Valid)
			require.True(t, smartCard.SantriID.Valid)
			require.Equal(t, santri.Name, smartCard.SantriName.String)
		}
	})

	t.Run("list all rfids should match employee", func(t *testing.T) {
		arg := ListSmartCardsParams{
			Q:            pgtype.Text{String: employee.Name[:3], Valid: true},
			OffsetNumber: 0,
			LimitNumber:  10,
			IsEmployee:   pgtype.Bool{Bool: true, Valid: true},
		}
		rfids, err := testStore.ListSmartCards(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, rfids)

		for _, smartCard := range rfids {
			require.NotEmpty(t, smartCard)
			require.Equal(t, randomEmployeeRfid.Uid, smartCard.Uid)
			require.Equal(t, randomEmployeeRfid.IsActive, smartCard.IsActive)
			require.Equal(t, randomEmployeeRfid.EmployeeID, smartCard.EmployeeID)

			//smartCard.santri_id should be null
			require.False(t, smartCard.SantriID.Valid)
			require.True(t, smartCard.EmployeeID.Valid)
			require.Equal(t, employee.Name, smartCard.EmployeeName.String)
		}
	})

}

func TestGetSmartCard(t *testing.T) {
	clearSmartCardTable(t)
	clearSantriTable(t)
	clearEmployeeTable(t)

	smartCard, _ := createRandomSmartCardWithSantri(t)
	getSmartCard, err := testStore.GetSmartCard(context.Background(), smartCard.Uid)
	require.NoError(t, err)
	require.NotEmpty(t, getSmartCard)
	require.Equal(t, smartCard.Uid, getSmartCard.Uid)
	require.Equal(t, smartCard.IsActive, getSmartCard.IsActive)

	require.Equal(t, smartCard.ID, getSmartCard.ID)
	require.Equal(t, smartCard.SantriID, getSmartCard.SantriID)
	require.Equal(t, smartCard.EmployeeID, getSmartCard.EmployeeID)

}

func TestUpdateSmartCard(t *testing.T) {
	clearSmartCardTable(t)
	clearSantriTable(t)
	clearEmployeeTable(t)

	smartCard, _ := createRandomSmartCardWithSantri(t)

	arg := UpdateSmartCardParams{
		ID:         smartCard.ID,
		Uid:        pgtype.Text{String: random.RandomString(12), Valid: true},
		IsActive:   pgtype.Bool{Bool: random.RandomBool(), Valid: true},
		SantriID:   pgtype.Int4{Int32: 0, Valid: false},
		EmployeeID: pgtype.Int4{Int32: 0, Valid: false},
	}
	updatedRfid, err := testStore.UpdateSmartCard(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedRfid)

	require.Equal(t, arg.ID, updatedRfid.ID)
	require.Equal(t, arg.Uid.String, updatedRfid.Uid)
	require.Equal(t, arg.IsActive.Bool, updatedRfid.IsActive)
	require.Equal(t, arg.SantriID.Int32, updatedRfid.SantriID.Int32)
	require.Equal(t, arg.EmployeeID.Int32, updatedRfid.EmployeeID.Int32)
}

func TestDeleteSmartCard(t *testing.T) {
	clearSmartCardTable(t)
	clearSantriTable(t)
	clearEmployeeTable(t)

	smartCard, _ := createRandomSmartCardWithSantri(t)

	deletedRfid, err := testStore.DeleteSmartCard(context.Background(), smartCard.ID)
	require.NoError(t, err)

	require.NotEmpty(t, deletedRfid)
}
