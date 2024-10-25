package db

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomEmployee(t *testing.T) Employee {
	occupation := createRandomEmployeeOccupation(t)
	arg := CreateEmployeeParams{
		Nip:          pgtype.Text{String: util.RandomString(18), Valid: true},
		Name:         util.RandomString(8),
		Gender:       GenderMale,
		Photo:        pgtype.Text{String: util.RandomString(12), Valid: true},
		OccupationID: occupation.ID,
	}
	employee, err := testQueries.CreateEmployee(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employee)

	require.Equal(t, arg.Nip.String, employee.Nip.String)
	require.Equal(t, arg.Name, employee.Name)
	require.Equal(t, arg.Gender, employee.Gender)
	require.Equal(t, arg.Photo.String, employee.Photo.String)
	require.Equal(t, arg.OccupationID, employee.OccupationID)

	return employee
}
func createRandomEmployeeWithUser(t *testing.T) (Employee, User) {
	occupation := createRandomEmployeeOccupation(t)
	user := createRandomUser(t)
	arg := CreateEmployeeParams{
		Nip:          pgtype.Text{String: util.RandomString(18), Valid: true},
		Name:         util.RandomString(8),
		Gender:       GenderMale,
		Photo:        pgtype.Text{String: util.RandomString(12), Valid: true},
		OccupationID: occupation.ID,
		UserID:       pgtype.Int4{Int32: user.ID, Valid: true},
	}
	employee, err := testQueries.CreateEmployee(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employee)

	require.Equal(t, arg.Nip.String, employee.Nip.String)
	require.Equal(t, arg.Name, employee.Name)
	require.Equal(t, arg.Gender, employee.Gender)
	require.Equal(t, arg.Photo.String, employee.Photo.String)
	require.Equal(t, arg.OccupationID, employee.OccupationID)

	return employee, user
}

func TestCreateEmployee(t *testing.T) {
	createRandomEmployee(t)
}

func TestQueryEmployeeWithQ(t *testing.T) {
	// Create test data with different names
	employee1 := createRandomEmployee(t)
	createRandomEmployee(t)
	createRandomEmployee(t)

	// Search for a specific parent name using `q`
	arg := QueryEmployeesAscParams{
		Q:            pgtype.Text{String: employee1.Name[:3], Valid: true},
		LimitNumber:  10,
		OffsetNumber: 0,
	}

	// Perform query
	employees, err := testQueries.QueryEmployeesAsc(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employees)

	// Verify that at least one result matches the queried name part
	found := false
	for _, employee := range employees {
		if employee.Name == employee1.Name {
			found = true
			break
		}
	}
	require.True(t, found, "Expected to find a employee matching the query")
}

func TestQueryEmployeeWithHasUser(t *testing.T) {

	_, user := createRandomEmployeeWithUser(t)
	createRandomEmployee(t)

	argWithUser := QueryEmployeesAscParams{
		HasUser:      1,
		LimitNumber:  10,
		OffsetNumber: 0,
	}
	employeessWithUser, err := testQueries.QueryEmployeesAsc(context.Background(), argWithUser)
	require.NoError(t, err)
	require.NotEmpty(t, employeessWithUser)

	for _, employee := range employeessWithUser {
		require.NotNil(t, employee.UserID, "Expected employee to have a user_id")
		if employee.UserID.Int32 == user.ID {
			require.Equal(t, user.Username.String, employee.UserUsername.String)
		}
	}

	// Query with `has_user = 0` (only parents without user_id)
	argWithoutUser := QueryEmployeesAscParams{
		HasUser:      0,
		LimitNumber:  10,
		OffsetNumber: 0,
	}
	employeesWithoutUser, err := testQueries.QueryEmployeesAsc(context.Background(), argWithoutUser)
	require.NoError(t, err)
	require.NotEmpty(t, employeesWithoutUser)

	for _, employee := range employeesWithoutUser {
		require.Zero(t, employee.UserID, "Expected employee to not have a user_id (0)")
	}

	// Query with `has_user = -1` (all parents)
	argAll := QueryEmployeesAscParams{
		HasUser:      -1,
		LimitNumber:  10,
		OffsetNumber: 0,
	}
	allEmployees, err := testQueries.QueryEmployeesAsc(context.Background(), argAll)
	require.NoError(t, err)
	require.NotEmpty(t, allEmployees)

	// Check that all parents are included
	hasUserCount := 0
	noUserCount := 0
	for _, employee := range allEmployees {
		if employee.UserID.Valid {
			hasUserCount++
		} else {
			noUserCount++
		}
	}
	require.GreaterOrEqual(t, len(allEmployees), 2, "Expected to retrieve all employees")
	require.GreaterOrEqual(t, hasUserCount, 1, "Expected to find parents with user_id")
	require.GreaterOrEqual(t, noUserCount, 1, "Expected to find parents without user_id")
}

func TestUpdateEmployee(t *testing.T) {
	employee1 := createRandomEmployee(t)

	// Update parent details
	newName := util.RandomString(8)
	newAddress := util.RandomString(50)
	newNoWa := util.RandomString(12)
	newPhoto := util.RandomString(12)

	arg := UpdateParentParams{
		ID:      employee1.ID,
		Name:    newName,
		Gender:  GenderMale,
		Address: newAddress,
		NoWa:    pgtype.Text{String: newNoWa, Valid: true},
		Photo:   pgtype.Text{String: newPhoto, Valid: true},
	}

	parent2, err := testQueries.UpdateParent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, parent2)

	require.Equal(t, employee1.ID, parent2.ID)
	require.Equal(t, newName, parent2.Name)
	require.Equal(t, newAddress, parent2.Address)
	require.Equal(t, employee1.Gender, parent2.Gender) // Gender should remain unchanged
	require.Equal(t, newNoWa, parent2.NoWa.String)
	require.Equal(t, newPhoto, parent2.Photo.String)
	require.Equal(t, employee1.UserID, parent2.UserID) // UserID should remain unchanged
}
