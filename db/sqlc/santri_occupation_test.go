package db

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearSantriOccupation(t *testing.T) {
	_, err := testQueries.db.Exec(context.Background(), `DELETE FROM "employee_occupation"`)
	require.NoError(t, err)
}

func createRandomSantriOccupation(t *testing.T) SantriOccupation {
	arg := CreateSantriOccupationParams{
		Name:        util.RandomString(8),
		Description: pgtype.Text{String: util.RandomString(50), Valid: true},
	}

	santriOccupation, err := testQueries.CreateSantriOccupation(context.Background(), arg)
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

func TestQuerySantriOccupation(t *testing.T) {
	clearSantriOccupation(t)
	createRandomSantriOccupation(t)
	createRandomSantriOccupation(t)
	createRandomSantriOccupation(t)

	santriOccupation, err := testQueries.QuerySantriOccupations(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, santriOccupation)

}

func TestUpdateSantriOccupation(t *testing.T) {
	clearSantriOccupation(t)
	santriOccupation := createRandomSantriOccupation(t)

	arg := UpdateSantriOccupationParams{
		ID:          santriOccupation.ID,
		Name:        util.RandomString(8),
		Description: pgtype.Text{String: util.RandomString(50), Valid: true},
	}

	santriOccupation, err := testQueries.UpdateSantriOccupation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, santriOccupation)

	require.Equal(t, arg.Name, santriOccupation.Name)
	require.Equal(t, arg.Description.String, santriOccupation.Description.String)
}

func TestDeleteSantriOccupation(t *testing.T) {
	clearSantriOccupation(t)
	santriOccupation := createRandomSantriOccupation(t)
	deletedEmployeeOccupation, err := testQueries.DeleteSantriOccupation(context.Background(), santriOccupation.ID)

	require.NoError(t, err)
	require.NotEmpty(t, deletedEmployeeOccupation)

	require.Equal(t, santriOccupation.ID, deletedEmployeeOccupation.ID)
}
