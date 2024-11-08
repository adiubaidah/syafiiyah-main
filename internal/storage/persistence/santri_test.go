package persistence

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearSantriTable(t *testing.T) {
	_, err := sqlStore.db.Exec(context.Background(), `DELETE FROM "santri"`)
	require.NoError(t, err)
}

func createRandomSantri(t *testing.T) Santri {

	arg := CreateSantriParams{
		Name:       random.RandomString(8),
		Nis:        pgtype.Text{String: random.RandomString(15), Valid: true},
		Gender:     GenderMale,
		Generation: int32(random.RandomInt(2010, 2030)),
		IsActive:   pgtype.Bool{Bool: random.RandomBool(), Valid: true},
		Photo:      pgtype.Text{String: random.RandomString(12), Valid: true},
	}
	santri, err := testStore.CreateSantri(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, santri)

	require.NotZero(t, santri.ID)
	require.Equal(t, arg.Name, santri.Name)
	require.Equal(t, arg.Nis.String, santri.Nis.String)
	require.Equal(t, arg.IsActive, santri.IsActive)
	require.Equal(t, arg.Gender, santri.Gender)
	require.Equal(t, arg.Generation, santri.Generation)
	require.Equal(t, arg.Photo.String, santri.Photo.String)

	return santri
}

func createRandomSantriWithParent(t *testing.T) (Santri, Parent) {
	parent := createRandomParent(t)
	arg := CreateSantriParams{
		Name:       random.RandomString(8),
		Nis:        pgtype.Text{String: random.RandomString(15), Valid: true},
		Gender:     GenderMale,
		Generation: int32(random.RandomInt(2010, 2030)),
		Photo:      pgtype.Text{String: random.RandomString(12), Valid: true},
		ParentID:   pgtype.Int4{Int32: parent.ID, Valid: true},
	}
	santri, err := testStore.CreateSantri(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, santri)

	require.Equal(t, arg.Name, santri.Name)
	require.Equal(t, arg.Nis.String, santri.Nis.String)
	require.Equal(t, arg.IsActive, santri.IsActive)
	require.Equal(t, arg.Gender, santri.Gender)
	require.Equal(t, arg.Generation, santri.Generation)
	require.Equal(t, arg.Photo.String, santri.Photo.String)

	return santri, parent
}
func TestCreateSantri(t *testing.T) {
	clearSantriPresenceTable(t)
	clearSantriPermissionTable(t)
	clearSantriTable(t)
	createRandomSantri(t)
}

func TestListSantri(t *testing.T) {
	clearSantriPresenceTable(t)
	clearSantriPermissionTable(t)
	clearSantriTable(t)
	randomSantri, _ := createRandomSantriWithParent(t)
	santris := []Santri{}
	for i := 0; i < 15; i++ {
		santris = append(santris, createRandomSantri(t))
	}

	t.Run("Run with List All", func(t *testing.T) {
		arg := ListSantriParams{
			LimitNumber:  15,
			OffsetNumber: 0,
		}

		allSantri, err := testStore.ListSantri(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, allSantri)
		require.Len(t, allSantri, 15)
	})

	t.Run("Run with List name", func(t *testing.T) {
		arg := ListSantriParams{
			Q:            pgtype.Text{String: randomSantri.Name[:3], Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
			OccupationID: pgtype.Int4{Int32: 0, Valid: false},
			Generation:   pgtype.Int4{Int32: 0, Valid: false},
			OrderBy:      NullSantriOrderBy{SantriOrderBy: SantriOrderByAscGeneration, Valid: true},
		}

		allSantri, err := testStore.ListSantri(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, allSantri)

		found := false
		for _, santri := range allSantri {
			if santri.Name == randomSantri.Name {
				found = true
				break
			}
		}
		require.True(t, found, "Expected to find a santri matching the List")
	})

	t.Run("list santri must contain active santri only", func(t *testing.T) {
		arg := ListSantriParams{
			Q:            pgtype.Text{String: "", Valid: false},
			IsActive:     pgtype.Bool{Bool: true, Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
			OccupationID: pgtype.Int4{Int32: 0, Valid: false},
			Generation:   pgtype.Int4{Int32: 0, Valid: false},
			OrderBy:      NullSantriOrderBy{SantriOrderBy: SantriOrderByAscGeneration, Valid: true},
		}
		result, err := testStore.ListSantri(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, result)

		for _, santri := range result {
			require.True(t, santri.IsActive.Bool)
		}
	})

	t.Run("list santri must contain not active santri only", func(t *testing.T) {
		arg := ListSantriParams{
			Q:            pgtype.Text{String: "", Valid: false},
			IsActive:     pgtype.Bool{Bool: false, Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
			OccupationID: pgtype.Int4{Int32: 0, Valid: false},
			Generation:   pgtype.Int4{Int32: 0, Valid: false},
			OrderBy:      NullSantriOrderBy{SantriOrderBy: SantriOrderByAscGeneration, Valid: true},
		}
		result, err := testStore.ListSantri(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, result)

		for _, santri := range result {
			require.False(t, santri.IsActive.Bool)
		}
	})

	t.Run("Run with List Nis", func(t *testing.T) {
		arg := ListSantriParams{
			Q:           pgtype.Text{String: randomSantri.Nis.String, Valid: true},
			LimitNumber: 10,

			OffsetNumber: 0,
		}

		allSantri, err := testStore.ListSantri(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, allSantri)

		found := false
		for _, santri := range allSantri {
			if santri.Name == randomSantri.Name {
				found = true
				break
			}
		}
		require.True(t, found, "Expected to find a santri matching the List")
	})
}

func TestListSantriPagination(t *testing.T) {
	clearSantriTable(t)
	for i := 0; i < 15; i++ {
		createRandomSantri(t)
	}

	testCases := []struct {
		name     string
		arg      ListSantriParams
		expected int
	}{
		{
			name: "Limit 5",
			arg: ListSantriParams{
				LimitNumber:  5,
				OffsetNumber: 0,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 5",
			arg: ListSantriParams{
				LimitNumber:  5,
				OffsetNumber: 5,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 10",
			arg: ListSantriParams{
				LimitNumber:  5,
				OffsetNumber: 10,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 10",
			arg: ListSantriParams{
				LimitNumber:  5,
				OffsetNumber: 15,
			},
			expected: 0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			allSantri, err := testStore.ListSantri(context.Background(), tt.arg)
			require.NoError(t, err)
			require.Len(t, allSantri, tt.expected)
		})
	}
}

func TestCountSantri(t *testing.T) {
	clearSantriTable(t)
	for i := 0; i < 15; i++ {
		createRandomSantri(t)
	}

	arg := CountSantriParams{
		Q:            pgtype.Text{String: "", Valid: false},
		OccupationID: pgtype.Int4{Int32: 0, Valid: false},
		Generation:   pgtype.Int4{Int32: 0, Valid: false},
		IsActive:     pgtype.Bool{Bool: true, Valid: false},
	}

	count, err := testStore.CountSantri(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, 15, int(count))
}

func TestGetSantri(t *testing.T) {
	clearSantriTable(t)
	santri := createRandomSantri(t)

	getSantri, err := testStore.GetSantri(context.Background(), santri.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getSantri)

	require.Equal(t, santri.ID, getSantri.ID)
	require.Equal(t, santri.Name, getSantri.Name)
	require.Equal(t, santri.Nis.String, getSantri.Nis.String)
	require.Equal(t, santri.IsActive.Bool, getSantri.IsActive.Bool)
	require.Equal(t, santri.Gender, getSantri.Gender)
	require.Equal(t, santri.Generation, getSantri.Generation)
	require.Equal(t, santri.Photo.String, getSantri.Photo.String)
	require.Equal(t, santri.ParentID.Int32, getSantri.ParentID.Int32)
	require.Equal(t, santri.OccupationID.Int32, getSantri.OccupationID.Int32)
}

func TestUpdateSantri(t *testing.T) {
	clearSantriTable(t)
	santri := createRandomSantri(t)

	arg := UpdateSantriParams{
		ID:         santri.ID,
		Name:       random.RandomString(8),
		Nis:        pgtype.Text{String: random.RandomString(15), Valid: true},
		IsActive:   false,
		Generation: int32(random.RandomInt(2010, 2030)),
		Photo:      pgtype.Text{String: random.RandomString(12), Valid: true},
	}

	updatedSantri, err := testStore.UpdateSantri(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSantri)

	require.Equal(t, arg.Name, updatedSantri.Name)
	require.Equal(t, arg.Nis.String, updatedSantri.Nis.String)
	require.Equal(t, arg.IsActive, updatedSantri.IsActive.Bool)
	require.Equal(t, arg.Generation, updatedSantri.Generation)
	require.Equal(t, arg.Photo.String, updatedSantri.Photo.String)
}

func TestDeleteSantri(t *testing.T) {
	clearSantriTable(t)
	santri := createRandomSantri(t)

	deletedSantri, err := testStore.DeleteSantri(context.Background(), santri.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedSantri)

	getSantri, err := testStore.GetSantri(context.Background(), santri.ID)
	require.Error(t, err)
	require.Empty(t, getSantri)

}
