package db

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearUserTable(t *testing.T) {
	_, err := testQueries.db.Exec(context.Background(), "DELETE FROM user")
	require.NoError(t, err)
}

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	roles := []UserRole{UserRoleSuperadmin, UserRoleAdmin, UserRoleEmployee, UserRoleParent}
	n := len(roles)
	arg := CreateUserParams{
		Username: pgtype.Text{String: util.RandomString(8), Valid: true},
		Role:     NullUserRole{UserRole: roles[util.RandomInt(0, int64(n-1))], Valid: true},
		Password: pgtype.Text{String: hashedPassword, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username.String, user.Username.String)
	require.Equal(t, arg.Role.UserRole, user.Role.UserRole)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
