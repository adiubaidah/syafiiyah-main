package repository

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearSantriOccupation(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), `DELETE FROM "employee_occupation"`)
	require.NoError(t, err)
}

func createRandomSantriOccupation(t *testing.T) SantriOccupation {
	arg := CreateSantriOccupationParams{
		Name:        random.RandomString(8),
		Description: pgtype.Text{String: random.RandomString(50), Valid: true},
	}

	santriOccupation, err := testStore.CreateSantriOccupation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, santriOccupation)

	require.Equal(t, arg.Name, santriOccupation.Name)
	require.Equal(t, arg.Description.String, santriOccupation.Description.String)
	return santriOccupation
}

func TestCreateSantriOccupation(t *testing.T) {
	clearSantriOccupation(t)
	createRandomSantriOccupation(t)
}

func TestListSantriOccupation(t *testing.T) {
	clearSantriOccupation(t)
	createRandomSantriOccupation(t)
	createRandomSantriOccupation(t)
	createRandomSantriOccupation(t)

	santriOccupation, err := testStore.ListSantriOccupations(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, santriOccupation)

}

func TestUpdateSantriOccupation(t *testing.T) {
	clearSantriOccupation(t)
	santriOccupation := createRandomSantriOccupation(t)

	arg := UpdateSantriOccupationParams{
		ID:          santriOccupation.ID,
		Name:        pgtype.Text{String: random.RandomString(8), Valid: true},
		Description: pgtype.Text{String: random.RandomString(50), Valid: true},
	}

	santriOccupation, err := testStore.UpdateSantriOccupation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, santriOccupation)

	require.Equal(t, arg.Name.String, santriOccupation.Name)
	require.Equal(t, arg.Description.String, santriOccupation.Description.String)
}

func TestDeleteSantriOccupation(t *testing.T) {
	clearSantriOccupation(t)
	santriOccupation := createRandomSantriOccupation(t)
	deletedEmployeeOccupation, err := testStore.DeleteSantriOccupation(context.Background(), santriOccupation.ID)

	require.NoError(t, err)
	require.NotEmpty(t, deletedEmployeeOccupation)

	require.Equal(t, santriOccupation.ID, deletedEmployeeOccupation.ID)
}
