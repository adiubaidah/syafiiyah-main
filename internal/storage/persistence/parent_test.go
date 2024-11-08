package persistence

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearParentTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), "DELETE FROM parent")
	require.NoError(t, err)
}

func createRandomParent(t *testing.T) Parent {
	arg := CreateParentParams{
		Name:           random.RandomString(8),
		Address:        random.RandomString(50),
		Gender:         GenderMale,
		WhatsappNumber: pgtype.Text{String: random.RandomString(12), Valid: true},
		Photo:          pgtype.Text{String: random.RandomString(12), Valid: true},
	}
	parent, err := testStore.CreateParent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, parent)

	require.Equal(t, arg.Name, parent.Name)
	require.Equal(t, arg.Address, parent.Address)
	require.Equal(t, arg.Gender, parent.Gender)
	require.Equal(t, arg.WhatsappNumber.String, parent.WhatsappNumber.String)
	require.Equal(t, arg.Photo.String, parent.Photo.String)

	return parent
}

func createRandomParentWithUser(t *testing.T) (Parent, User) {
	user := createRandomUser(t, UserRoleParent)
	arg := CreateParentParams{
		Name:           random.RandomString(8),
		Address:        random.RandomString(50),
		Gender:         GenderMale,
		WhatsappNumber: pgtype.Text{String: random.RandomString(12), Valid: true},
		Photo:          pgtype.Text{String: random.RandomString(12), Valid: true},
		UserID:         pgtype.Int4{Int32: user.ID, Valid: true},
	}
	parent, err := testStore.CreateParent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, parent)

	require.Equal(t, arg.Name, parent.Name)
	require.Equal(t, arg.Address, parent.Address)
	require.Equal(t, arg.Gender, parent.Gender)
	require.Equal(t, arg.WhatsappNumber.String, parent.WhatsappNumber.String)
	require.Equal(t, arg.Photo.String, parent.Photo.String)
	require.Equal(t, arg.UserID.Int32, parent.UserID.Int32)
	return parent, user
}

func TestCreateParent(t *testing.T) {
	createRandomParent(t)
}

func TestListParentsWithQ(t *testing.T) {
	clearParentTable(t)
	// Create test data with different names
	parent1 := createRandomParent(t)
	createRandomParent(t)
	createRandomParent(t)

	t.Run("List by partial parent name", func(t *testing.T) {
		// Search for a specific parent name using `q`
		arg := ListParentsAscParams{
			Q:            pgtype.Text{String: parent1.Name[:3], Valid: true}, // Partially match the first 3 characters of name
			LimitNumber:  10,
			OffsetNumber: 0,
		}

		// Perform List
		parents, err := testStore.ListParentsAsc(context.Background(), arg)
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
		require.True(t, found, "Expected to find a parent matching the List")
	})
}

func TestListParentWithHasUser(t *testing.T) {
	clearParentTable(t)
	// Create test data with and without user IDs
	_, user := createRandomParentWithUser(t)
	createRandomParent(t)

	t.Run("List parents with user ID", func(t *testing.T) {
		// List with `has_user = 1` (only parents with user_id)
		arg := ListParentsAscParams{
			HasUser:      1,
			LimitNumber:  10,
			OffsetNumber: 0,
		}

		parents, err := testStore.ListParentsAsc(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, parents)

		for _, parent := range parents {
			require.NotNil(t, parent.UserID, "Expected parent to have a user_id")
			if parent.UserID.Int32 == user.ID {
				require.Equal(t, user.Username.String, parent.UserUsername.String)
			}
		}
	})

	t.Run("List parents without user ID", func(t *testing.T) {
		// List with `has_user = 0` (only parents without user_id)
		arg := ListParentsAscParams{
			HasUser:      0,
			LimitNumber:  10,
			OffsetNumber: 0,
		}

		parents, err := testStore.ListParentsAsc(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, parents)

		for _, parent := range parents {
			require.Zero(t, parent.UserID, "Expected parent to not have a user_id (0)")
		}
	})

	t.Run("List all parents regardless of user ID", func(t *testing.T) {
		// List with `has_user = -1` (all parents)
		arg := ListParentsAscParams{
			HasUser:      -1,
			LimitNumber:  10,
			OffsetNumber: 0,
		}

		parents, err := testStore.ListParentsAsc(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, parents)

		// Check that all parents are included
		hasUserCount := 0
		noUserCount := 0
		for _, parent := range parents {
			if parent.UserID.Valid {
				hasUserCount++
			} else {
				noUserCount++
			}
		}
		require.GreaterOrEqual(t, len(parents), 2, "Expected to retrieve all parents")
		require.GreaterOrEqual(t, hasUserCount, 1, "Expected to find parents with user_id")
		require.GreaterOrEqual(t, noUserCount, 1, "Expected to find parents without user_id")
	})
}

func TestListParentPagination(t *testing.T) {
	clearParentTable(t)
	for i := 0; i < 10; i++ {
		createRandomParent(t)
	}

	testCases := []struct {
		name     string
		arg      ListParentsAscParams
		expected int
	}{
		{
			name: "Limit 5",
			arg: ListParentsAscParams{
				LimitNumber:  5,
				OffsetNumber: 0,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 5",
			arg: ListParentsAscParams{
				LimitNumber:  5,
				OffsetNumber: 5,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 10",
			arg: ListParentsAscParams{
				LimitNumber:  5,
				OffsetNumber: 10,
			},
			expected: 0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			employees, err := testStore.ListParentsAsc(context.Background(), tt.arg)
			require.NoError(t, err)
			require.Len(t, employees, tt.expected)
		})
	}
}

func TestCountParents(t *testing.T) {
	clearParentTable(t)
	// Create test data
	createRandomParent(t)
	createRandomParent(t)
	createRandomParent(t)

	// Count parents
	arg := CountParentsParams{
		Q:       pgtype.Text{Valid: false},
		HasUser: -1,
	}
	count, err := testStore.CountParents(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, count)

	require.Greater(t, count, int64(2))
}

func TestUpdateParent(t *testing.T) {
	clearParentTable(t)
	parent1 := createRandomParent(t)

	// Update parent details
	newName := random.RandomString(8)
	newAddress := random.RandomString(50)
	newNoWa := random.RandomString(12)
	newPhoto := random.RandomString(12)

	arg := UpdateParentParams{
		ID:             parent1.ID,
		Name:           newName,
		Gender:         GenderMale,
		Address:        newAddress,
		WhatsappNumber: pgtype.Text{String: newNoWa, Valid: true},
		Photo:          pgtype.Text{String: newPhoto, Valid: true},
	}

	parent2, err := testStore.UpdateParent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, parent2)

	require.Equal(t, parent1.ID, parent2.ID)
	require.Equal(t, newName, parent2.Name)
	require.Equal(t, newAddress, parent2.Address)
	require.Equal(t, parent1.Gender, parent2.Gender) // Gender should remain unchanged
	require.Equal(t, newNoWa, parent2.WhatsappNumber.String)
	require.Equal(t, newPhoto, parent2.Photo.String)
	require.Equal(t, parent1.UserID, parent2.UserID) // UserID should remain unchanged
}

func TestDeleteParent(t *testing.T) {
	clearParentTable(t)
	parent := createRandomParent(t)

	deletedParent, err := testStore.DeleteParent(context.Background(), parent.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedParent)

	require.Equal(t, parent.ID, deletedParent.ID)

	parent2, err := testStore.GetParent(context.Background(), parent.ID)
	require.Error(t, err)
	require.Empty(t, parent2)
}
