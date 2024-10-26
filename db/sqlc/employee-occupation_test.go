package db

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearEmployeeOccupationTable(t *testing.T) {
	_, err := testQueries.db.Exec(context.Background(), `DELETE FROM "employee_occupation"`)
	require.NoError(t, err)
}

func createRandomEmployeeOccupation(t *testing.T) EmployeeOccupation {
	arg := CreateEmployeeOccupationParams{
		Name:        util.RandomString(8),
		Description: pgtype.Text{String: util.RandomString(50), Valid: true},
	}

	employeeOccupation, err := testQueries.CreateEmployeeOccupation(context.Background(), arg)
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

func TestQueryEmployeeOccupation(t *testing.T) {
	clearEmployeeOccupationTable(t)
	createRandomEmployeeOccupation(t)
	createRandomEmployeeOccupation(t)
	createRandomEmployeeOccupation(t)

	employeeOccupations, err := testQueries.QueryEmployeeOccupations(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, employeeOccupations)
}

func TestUpdateEmployeeOccupation(t *testing.T) {
	clearEmployeeOccupationTable(t)
	employeeOccupation := createRandomEmployeeOccupation(t)

	arg := UpdateEmployeeOccupationParams{
		ID:          employeeOccupation.ID,
		Name:        util.RandomString(8),
		Description: pgtype.Text{String: util.RandomString(50), Valid: true},
	}

	employeeOccupation, err := testQueries.UpdateEmployeeOccupation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employeeOccupation)

	require.Equal(t, arg.Name, employeeOccupation.Name)
	require.Equal(t, arg.Description.String, employeeOccupation.Description.String)
}

func TestDeleteEmployeeOccupation(t *testing.T) {
	clearEmployeeOccupationTable(t)
	employeeOccupation := createRandomEmployeeOccupation(t)
	deletedEmployeeOccupation, err := testQueries.DeleteEmployeeOccupation(context.Background(), employeeOccupation.ID)

	require.NoError(t, err)
	require.NotEmpty(t, deletedEmployeeOccupation)

	require.Equal(t, employeeOccupation.ID, deletedEmployeeOccupation.ID)
}
