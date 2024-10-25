package db

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearParentTable(t *testing.T) {
	clearUserTable(t)
	_, err := testQueries.db.Exec(context.Background(), "DELETE FROM parent")
	require.NoError(t, err)
}

func createRandomParent(t *testing.T) Parent {
	arg := CreateParentParams{
		Name:    util.RandomString(8),
		Address: util.RandomString(50),
		Gender:  GenderMale,
		NoWa:    pgtype.Text{String: util.RandomString(12), Valid: true},
		Photo:   pgtype.Text{String: util.RandomString(12), Valid: true},
	}
	parent, err := testQueries.CreateParent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, parent)

	require.Equal(t, arg.Name, parent.Name)
	require.Equal(t, arg.Address, parent.Address)
	require.Equal(t, arg.Gender, parent.Gender)
	require.Equal(t, arg.NoWa.String, parent.NoWa.String)
	require.Equal(t, arg.Photo.String, parent.Photo.String)

	return parent
}

func createRandomParentWithUser(t *testing.T) (Parent, User) {
	user := createRandomUser(t)
	arg := CreateParentParams{
		Name:    util.RandomString(8),
		Address: util.RandomString(50),
		Gender:  GenderMale,
		NoWa:    pgtype.Text{String: util.RandomString(12), Valid: true},
		Photo:   pgtype.Text{String: util.RandomString(12), Valid: true},
		UserID:  pgtype.Int4{Int32: user.ID, Valid: true},
	}
	parent, err := testQueries.CreateParent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, parent)

	require.Equal(t, arg.Name, parent.Name)
	require.Equal(t, arg.Address, parent.Address)
	require.Equal(t, arg.Gender, parent.Gender)
	require.Equal(t, arg.NoWa.String, parent.NoWa.String)
	require.Equal(t, arg.Photo.String, parent.Photo.String)
	require.Equal(t, arg.UserID.Int32, parent.UserID.Int32)
	return parent, user
}

func TestCreateParent(t *testing.T) {
	createRandomParent(t)
}

func TestQueryParentsWithQ(t *testing.T) {
	// Create test data with different names
	parent1 := createRandomParent(t)
	createRandomParent(t)
	createRandomParent(t)

	// Search for a specific parent name using `q`
	arg := QueryParentsAscParams{
		Q:            pgtype.Text{String: parent1.Name[:3], Valid: true}, // Partially match the first 3 characters of name                                // Ignore `has_user` filter
		LimitNumber:  10,
		OffsetNumber: 0,
	}

	// Perform query
	parents, err := testQueries.QueryParentsAsc(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, parents)

	// Verify that at least one result matches the queried name part
	found := false
	for _, parent := range parents {
		if parent.Name == parent1.Name {
			found = true
			break
		}
	}
	require.True(t, found, "Expected to find a parent matching the query")
}

func TestQueryParentWithHasUser(t *testing.T) {
	// Create test data with and without user IDs
	_, user := createRandomParentWithUser(t)
	createRandomParent(t)

	// Query with `has_user = 1` (only parents with user_id)
	argWithUser := QueryParentsAscParams{
		HasUser:      1,
		LimitNumber:  10,
		OffsetNumber: 0,
	}
	parentsWithUser, err := testQueries.QueryParentsAsc(context.Background(), argWithUser)
	require.NoError(t, err)
	require.NotEmpty(t, parentsWithUser)

	for _, parent := range parentsWithUser {
		require.NotNil(t, parent.UserID, "Expected parent to have a user_id")
		if parent.UserID.Int32 == user.ID {
			require.Equal(t, user.Username.String, parent.UserUsername.String)
		}
	}

	// Query with `has_user = 0` (only parents without user_id)
	argWithoutUser := QueryParentsAscParams{
		HasUser:      0,
		LimitNumber:  10,
		OffsetNumber: 0,
	}
	parentsWithoutUser, err := testQueries.QueryParentsAsc(context.Background(), argWithoutUser)
	require.NoError(t, err)
	require.NotEmpty(t, parentsWithoutUser)

	for _, parent := range parentsWithoutUser {
		require.Zero(t, parent.UserID, "Expected parent to not have a user_id (0)")
	}

	// Query with `has_user = -1` (all parents)
	argAll := QueryParentsAscParams{
		HasUser:      -1,
		LimitNumber:  10,
		OffsetNumber: 0,
	}
	allParents, err := testQueries.QueryParentsAsc(context.Background(), argAll)
	require.NoError(t, err)
	require.NotEmpty(t, allParents)

	// Check that all parents are included
	hasUserCount := 0
	noUserCount := 0
	for _, parent := range allParents {
		if parent.UserID.Valid {
			hasUserCount++
		} else {
			noUserCount++
		}
	}
	require.GreaterOrEqual(t, len(allParents), 2, "Expected to retrieve all parents")
	require.GreaterOrEqual(t, hasUserCount, 1, "Expected to find parents with user_id")
	require.GreaterOrEqual(t, noUserCount, 1, "Expected to find parents without user_id")
}
func TestUpdateParent(t *testing.T) {
	parent1 := createRandomParent(t)

	// Update parent details
	newName := util.RandomString(8)
	newAddress := util.RandomString(50)
	newNoWa := util.RandomString(12)
	newPhoto := util.RandomString(12)

	arg := UpdateParentParams{
		ID:      parent1.ID,
		Name:    newName,
		Gender:  GenderMale,
		Address: newAddress,
		NoWa:    pgtype.Text{String: newNoWa, Valid: true},
		Photo:   pgtype.Text{String: newPhoto, Valid: true},
	}

	parent2, err := testQueries.UpdateParent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, parent2)

	require.Equal(t, parent1.ID, parent2.ID)
	require.Equal(t, newName, parent2.Name)
	require.Equal(t, newAddress, parent2.Address)
	require.Equal(t, parent1.Gender, parent2.Gender) // Gender should remain unchanged
	require.Equal(t, newNoWa, parent2.NoWa.String)
	require.Equal(t, newPhoto, parent2.Photo.String)
	require.Equal(t, parent1.UserID, parent2.UserID) // UserID should remain unchanged
}

func TestDeleteParent(t *testing.T) {
	parent := createRandomParent(t)

	deletedParent, err := testQueries.DeleteParent(context.Background(), parent.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedParent)

	require.Equal(t, parent.ID, deletedParent.ID)

	parent2, err := testQueries.GetParent(context.Background(), parent.ID)
	require.Error(t, err)
	require.Empty(t, parent2)
}
