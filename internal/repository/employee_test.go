package repository

import (
	"context"
	"testing"

	"github.com/adiubaidah/syafiiyah-main/pkg/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearEmployeeTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), `DELETE FROM "employee"`)
	require.NoError(t, err)
}

func createRandomEmployee(t *testing.T) Employee {
	occupation := createRandomEmployeeOccupation(t)
	arg := CreateEmployeeParams{
		Nip:          pgtype.Text{String: random.RandomString(18), Valid: true},
		Name:         random.RandomString(8),
		Gender:       GenderTypeMale,
		Photo:        pgtype.Text{String: random.RandomString(12), Valid: true},
		OccupationID: occupation.ID,
	}
	employee, err := testStore.CreateEmployee(context.Background(), arg)
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
	user := createRandomUser(t, RoleTypeEmployee)
	arg := CreateEmployeeParams{
		Nip:          pgtype.Text{String: random.RandomString(18), Valid: true},
		Name:         random.RandomString(8),
		Gender:       GenderTypeMale,
		Photo:        pgtype.Text{String: random.RandomString(12), Valid: true},
		OccupationID: occupation.ID,
		UserID:       pgtype.Int4{Int32: user.ID, Valid: true},
	}
	employee, err := testStore.CreateEmployee(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employee)

	require.Equal(t, arg.Nip.String, employee.Nip.String)
	require.Equal(t, arg.Name, employee.Name)
	require.Equal(t, arg.Gender, employee.Gender)
	require.Equal(t, arg.Photo.String, employee.Photo.String)
	require.Equal(t, arg.OccupationID, employee.OccupationID)
	require.NotZero(t, employee.UserID.Int32)
	require.Equal(t, arg.UserID.Int32, employee.UserID.Int32)

	return employee, user
}

func TestCreateEmployee(t *testing.T) {
	createRandomEmployee(t)
}

func TestListEmployeeWithQ(t *testing.T) {
	clearEmployeeTable(t)
	employee1 := createRandomEmployee(t)
	createRandomEmployee(t)
	createRandomEmployee(t)

	count := 0

	t.Run("List with `q` (search by name)", func(t *testing.T) {
		arg := ListEmployeesParams{
			Q:            pgtype.Text{String: employee1.Name[:3], Valid: true},
			HasUser:      pgtype.Bool{Valid: false},
			LimitNumber:  10,
			OffsetNumber: 0,
			OrderBy:      NullEmployeeOrderBy{Valid: false},
			OccupationID: pgtype.Int4{Valid: false},
		}

		employees, err := testStore.ListEmployees(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, employees)

		found := false
		for _, employee := range employees {
			if employee.Name == employee1.Name {
				found = true
			}
			count++
		}
		require.True(t, found, "Expected to find a employee matching the List")
	})

	t.Run("count employee", func(t *testing.T) {
		result, err := testStore.CountEmployees(context.Background(), CountEmployeesParams{
			Q:            pgtype.Text{String: employee1.Name[:3], Valid: true},
			HasUser:      pgtype.Bool{Valid: false},
			OccupationID: pgtype.Int4{Valid: false},
		})
		require.NoError(t, err)
		require.Equal(t, count, int(result))
	})
}

func TestListEmployeeWithHasUser(t *testing.T) {
	clearEmployeeTable(t)
	_, user := createRandomEmployeeWithUser(t)
	createRandomEmployee(t)

	t.Run("List with `has_user = 1` (only employees with user_id)", func(t *testing.T) {
		arg := ListEmployeesParams{
			HasUser:      pgtype.Bool{Bool: true, Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
		}
		employeessWithUser, err := testStore.ListEmployees(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, employeessWithUser)

		for _, employee := range employeessWithUser {
			require.NotNil(t, employee.UserID, "Expected employee to have a user_id")
			if employee.UserID.Int32 == user.ID {
				require.Equal(t, user.Username.String, employee.Username.String)
			}
		}
	})

	t.Run("List with `has_user = 0` (only parents without user_id)", func(t *testing.T) {
		arg := ListEmployeesParams{
			HasUser:      pgtype.Bool{Bool: false, Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
		}
		employeesWithoutUser, err := testStore.ListEmployees(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, employeesWithoutUser)

		for _, employee := range employeesWithoutUser {
			require.Zero(t, employee.UserID, "Expected employee to not have a user_id (0)")
		}
	})

	t.Run("List with `has_user = -1` (all employee)", func(t *testing.T) {
		arg := ListEmployeesParams{
			HasUser:      pgtype.Bool{Valid: false},
			LimitNumber:  10,
			OffsetNumber: 0,
		}
		allEmployees, err := testStore.ListEmployees(context.Background(), arg)
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
		require.GreaterOrEqual(t, hasUserCount, 1, "Expected to find parents that has user_id")
		require.GreaterOrEqual(t, noUserCount, 1, "Expected to find parents without user_id")
	})
}

func TestListEmployeePagination(t *testing.T) {
	clearEmployeeTable(t)
	for i := 0; i < 10; i++ {
		createRandomEmployee(t)
	}

	testCases := []struct {
		name     string
		arg      ListEmployeesParams
		expected int
	}{
		{
			name: "Limit 5",
			arg: ListEmployeesParams{
				LimitNumber:  5,
				OffsetNumber: 0,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 5",
			arg: ListEmployeesParams{
				LimitNumber:  5,
				OffsetNumber: 5,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 10",
			arg: ListEmployeesParams{
				LimitNumber:  5,
				OffsetNumber: 10,
			},
			expected: 0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			employees, err := testStore.ListEmployees(context.Background(), tt.arg)
			require.NoError(t, err)
			require.Len(t, employees, tt.expected)
		})
	}
}

func TestGetEmployee(t *testing.T) {
	employee1 := createRandomEmployee(t)

	employee2, err := testStore.GetEmployeeByID(context.Background(), employee1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, employee2)

	require.Equal(t, employee1.ID, employee2.ID)
	require.Equal(t, employee1.Nip.String, employee2.Nip.String)
	require.Equal(t, employee1.Name, employee2.Name)
	require.Equal(t, employee1.Gender, employee2.Gender)
	require.Equal(t, employee1.Photo.String, employee2.Photo.String)
	require.Equal(t, employee1.OccupationID, employee2.OccupationID)
	require.Equal(t, employee1.UserID, employee2.UserID)
}

func TestUpdateEmployee(t *testing.T) {
	clearEmployeeTable(t)
	employee1 := createRandomEmployee(t)

	// Update parent details
	newPhoto := random.RandomString(12)

	arg := UpdateEmployeeParams{
		ID:           employee1.ID,
		Nip:          pgtype.Text{String: random.RandomString(18), Valid: true},
		Name:         pgtype.Text{String: random.RandomString(8), Valid: true},
		Gender:       NullGenderType{Valid: false},
		Photo:        pgtype.Text{String: newPhoto, Valid: true},
		OccupationID: pgtype.Int4{Int32: employee1.OccupationID, Valid: true},
		UserID:       pgtype.Int4{Int32: employee1.UserID.Int32, Valid: employee1.UserID.Valid},
	}

	employeeUpdated, err := testStore.UpdateEmployee(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employeeUpdated)

	require.Equal(t, employee1.ID, employeeUpdated.ID)
	require.Equal(t, arg.Name.String, employeeUpdated.Name)
	require.Equal(t, arg.Nip.String, employeeUpdated.Nip.String)
	require.Equal(t, employee1.OccupationID, employeeUpdated.OccupationID)
	require.Equal(t, employee1.Gender, employeeUpdated.Gender) // Gender should remain unchanged
	require.Equal(t, newPhoto, employeeUpdated.Photo.String)
	require.Equal(t, employee1.UserID, employeeUpdated.UserID) // UserID should remain unchanged
}

func TestUpdateBindEmployeeAndUser(t *testing.T) {
	clearEmployeeTable(t)
	clearUserTable(t)
	employee, _ := createRandomEmployeeWithUser(t)

	t.Run("Bind user to two employees and must failed", func(t *testing.T) {
		employee2 := createRandomEmployee(t)
		arg := UpdateEmployeeParams{
			ID:           employee2.ID,
			Nip:          pgtype.Text{String: employee2.Nip.String, Valid: true},
			Name:         pgtype.Text{String: employee2.Name, Valid: true},
			Gender:       NullGenderType{Valid: false},
			Photo:        employee2.Photo,
			OccupationID: pgtype.Int4{Int32: employee2.OccupationID, Valid: true},
			UserID:       pgtype.Int4{Int32: employee.UserID.Int32, Valid: true},
		}

		_, err := testStore.UpdateEmployee(context.Background(), arg)
		require.Error(t, err)
	})

	t.Run("Unbind employee from user", func(t *testing.T) {
		arg := UpdateEmployeeParams{
			ID:           employee.ID,
			Nip:          pgtype.Text{String: employee.Nip.String, Valid: true},
			Name:         pgtype.Text{String: employee.Name, Valid: true},
			Gender:       NullGenderType{Valid: false},
			Photo:        employee.Photo,
			OccupationID: pgtype.Int4{Int32: employee.OccupationID, Valid: true},
			UserID:       pgtype.Int4{Valid: false},
		}

		employeeUpdated, err := testStore.UpdateEmployee(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, employeeUpdated)
		require.Zero(t, employeeUpdated.UserID.Int32, "Expected employee to be unbinded from user")
	})
}

func TestDeleteEmployee(t *testing.T) {
	clearEmployeeTable(t)
	employee := createRandomEmployee(t)
	deletedEmployee, err := testStore.DeleteEmployee(context.Background(), employee.ID)

	require.NoError(t, err)
	require.NotEmpty(t, deletedEmployee)

	// Check that the deleted employee is the same as the one we created
	assert.Equal(t, employee.ID, deletedEmployee.ID)
	assert.Equal(t, employee.Nip.String, deletedEmployee.Nip.String)
	assert.Equal(t, employee.Name, deletedEmployee.Name)
	assert.Equal(t, employee.Gender, deletedEmployee.Gender)
	assert.Equal(t, employee.Photo.String, deletedEmployee.Photo.String)
	assert.Equal(t, employee.OccupationID, deletedEmployee.OccupationID)
	assert.Equal(t, employee.UserID, deletedEmployee.UserID)
}
