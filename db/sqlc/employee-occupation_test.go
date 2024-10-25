package db

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

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
	createRandomEmployeeOccupation(t)
}

func TestQueryEmployeeOccupation(t *testing.T) {
	createRandomEmployeeOccupation(t)
	createRandomEmployeeOccupation(t)
	createRandomEmployeeOccupation(t)

	employeeOccupations, err := testQueries.QueryEmployeeOccupations(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, employeeOccupations)
}

func TestUpdateEmployeeOccupation(t *testing.T) {
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
	employeeOccupation := createRandomEmployeeOccupation(t)
	deletedEmployeeOccupation, err := testQueries.DeleteEmployeeOccupation(context.Background(), employeeOccupation.ID)

	require.NoError(t, err)
	require.NotEmpty(t, deletedEmployeeOccupation)

	require.Equal(t, employeeOccupation.ID, deletedEmployeeOccupation.ID)
}

// func TestDeleteEmployeeOccupation(t *testing.T) {
// 	employeeOccupation := createRandomEmployeeOccupation(t)

// 	err := testQueries.DeleteEmployeeOccupation(context.Background(), employeeOccupation.ID)
// 	require.NoError(t, err)

// 	employeeOccupation, err = testQueries.GetEmployeeOccupation(context.Background(), employeeOccupation.ID)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrNotFound.Error())
// 	require.Empty(t, employeeOccupation)
// }
