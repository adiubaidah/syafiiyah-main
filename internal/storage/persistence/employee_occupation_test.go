package persistence

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearEmployeeOccupationTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), `DELETE FROM "employee_occupation"`)
	require.NoError(t, err)
}

func createRandomEmployeeOccupation(t *testing.T) EmployeeOccupation {
	arg := CreateEmployeeOccupationParams{
		Name:        random.RandomString(8),
		Description: pgtype.Text{String: random.RandomString(50), Valid: true},
	}

	employeeOccupation, err := testStore.CreateEmployeeOccupation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employeeOccupation)

	require.Equal(t, arg.Name, employeeOccupation.Name)
	require.Equal(t, arg.Description.String, employeeOccupation.Description.String)
	return employeeOccupation
}

func TestCreateEmployeeOccupation(t *testing.T) {
	clearEmployeeOccupationTable(t)
	createRandomEmployeeOccupation(t)
}

func TestListEmployeeOccupation(t *testing.T) {
	clearEmployeeOccupationTable(t)
	createRandomEmployeeOccupation(t)
	createRandomEmployeeOccupation(t)
	createRandomEmployeeOccupation(t)

	employeeOccupations, err := testStore.ListEmployeeOccupations(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, employeeOccupations)
}

func TestUpdateEmployeeOccupation(t *testing.T) {
	clearEmployeeOccupationTable(t)
	employeeOccupation := createRandomEmployeeOccupation(t)

	arg := UpdateEmployeeOccupationParams{
		ID:          employeeOccupation.ID,
		Name:        random.RandomString(8),
		Description: pgtype.Text{String: random.RandomString(50), Valid: true},
	}

	employeeOccupation, err := testStore.UpdateEmployeeOccupation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employeeOccupation)

	require.Equal(t, arg.Name, employeeOccupation.Name)
	require.Equal(t, arg.Description.String, employeeOccupation.Description.String)
}

func TestDeleteEmployeeOccupation(t *testing.T) {
	clearEmployeeOccupationTable(t)
	employeeOccupation := createRandomEmployeeOccupation(t)
	deletedEmployeeOccupation, err := testStore.DeleteEmployeeOccupation(context.Background(), employeeOccupation.ID)

	require.NoError(t, err)
	require.NotEmpty(t, deletedEmployeeOccupation)

	require.Equal(t, employeeOccupation.ID, deletedEmployeeOccupation.ID)
}
