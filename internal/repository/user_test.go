package repository

import (
	"context"
	"testing"

	"github.com/adiubaidah/syafiiyah-main/pkg/random"
	"github.com/adiubaidah/syafiiyah-main/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearUserTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), `DELETE FROM "user"`)
	require.NoError(t, err)
}

func createRandomUser(t *testing.T, role RoleType) User {
	hashedPassword, err := util.HashPassword(random.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username: random.RandomString(8),
		Email:    random.RandomEmail(),
		Role:     role,
		Password: hashedPassword,
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username.String)
	require.Equal(t, arg.Email, user.Email.String)
	require.Equal(t, arg.Role, user.Role.RoleType)

	return user
}

func TestCreateUser(t *testing.T) {
	clearParentTable(t)
	clearEmployeeTable(t)
	clearUserTable(t)
	roles := []RoleType{RoleTypeSuperadmin, RoleTypeAdmin, RoleTypeEmployee, RoleTypeParent}
	n := len(roles)
	createRandomUser(t, roles[random.RandomInt(0, int64(n-1))])
}

func TestListUsersOwner(t *testing.T) {
	clearEmployeeTable(t)
	clearParentTable(t)
	clearUserTable(t)

	parent1, parentUser1 := createRandomParentWithUser(t)
	employee1, employeeUser1 := createRandomEmployeeWithUser(t)

	t.Run("Should return users matching List", func(t *testing.T) {
		arg := ListUserParams{
			Q:            pgtype.Text{String: parentUser1.Username.String, Valid: true},
			Role:         NullRoleType{Valid: false},
			LimitNumber:  10,
			OffsetNumber: 0,
			HasOwner:     pgtype.Bool{Valid: false},
		}

		users, err := testStore.ListUsers(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, users)
		require.Len(t, users, 1)
	})

	t.Run("list user should match role", func(t *testing.T) {
		arg := ListUserParams{
			Q:            pgtype.Text{String: "", Valid: false},
			Role:         NullRoleType{Valid: true, RoleType: RoleTypeParent},
			LimitNumber:  10,
			OffsetNumber: 0,
			HasOwner:     pgtype.Bool{Valid: false},
		}
		users, err := testStore.ListUsers(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, users)

		for _, user := range users {
			require.Equal(t, RoleTypeParent, user.Role)
		}
	})

	t.Run("should return users with owner", func(t *testing.T) {

		arg := ListUserParams{
			Q:            pgtype.Text{String: "", Valid: false},
			Role:         NullRoleType{Valid: false},
			LimitNumber:  10,
			OffsetNumber: 0,
			HasOwner:     pgtype.Bool{Valid: true, Bool: true},
		}
		users, err := testStore.ListUsers(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, users)

		for _, user := range users {
			require.NotZero(t, user.IDOwner.Int32)
			if user.ID == parentUser1.ID {
				require.Equal(t, parent1.ID, user.IDOwner.Int32)
			}
			if user.ID == employeeUser1.ID {
				require.Equal(t, employee1.ID, user.IDOwner.Int32)
			}
		}
	})
}

func TestListUserPagination(t *testing.T) {
	clearUserTable(t)
	for i := 0; i < 10; i++ {
		createRandomUser(t, RoleTypeSuperadmin)
	}

	testCases := []struct {
		name     string
		arg      ListUserParams
		expected int
	}{
		{
			name: "Limit 5",
			arg: ListUserParams{
				LimitNumber:  5,
				OffsetNumber: 0,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 5",
			arg: ListUserParams{
				LimitNumber:  5,
				OffsetNumber: 5,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 10",
			arg: ListUserParams{
				LimitNumber:  5,
				OffsetNumber: 10,
			},
			expected: 0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			users, err := testStore.ListUsers(context.Background(), tt.arg)
			require.NoError(t, err)
			require.Len(t, users, tt.expected)
		})
	}
}

func TestCountUsers(t *testing.T) {
	clearUserTable(t)
	for i := 0; i < 5; i++ {
		createRandomUser(t, RoleTypeSuperadmin)
	}

	count, err := testStore.CountUsers(context.Background(), CountUsersParams{})
	require.NoError(t, err)
	require.Equal(t, int32(5), int32(count))
}

func TestUpdateUser(t *testing.T) {
	clearUserTable(t)
	user1 := createRandomUser(t, RoleTypeSuperadmin)

	arg := UpdateUserParams{
		Username: pgtype.Text{String: random.RandomString(8), Valid: true},
		Role:     NullRoleType{Valid: true, RoleType: RoleTypeAdmin},
		ID:       user1.ID,
	}

	user, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username.String, user.Username.String)
	require.Equal(t, arg.Role.RoleType, user.Role.RoleType)
}

func TestDeleteUser(t *testing.T) {
	clearUserTable(t)
	user1 := createRandomUser(t, RoleTypeSuperadmin)

	userDeleted, err := testStore.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, userDeleted)

	require.Equal(t, user1.ID, userDeleted.ID)
	require.Equal(t, user1.Username.String, userDeleted.Username.String)
	require.Equal(t, user1.Role.RoleType, userDeleted.Role.RoleType)
}
