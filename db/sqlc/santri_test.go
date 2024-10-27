package db

import (
	"context"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func clearSantriTable(t *testing.T) {
	_, err := testQueries.db.Exec(context.Background(), `DELETE FROM "santri"`)
	require.NoError(t, err)
}

func createRandomSantri(t *testing.T) Santri {

	arg := CreateSantriParams{
		Name:       util.RandomString(8),
		Nis:        pgtype.Text{String: util.RandomString(15), Valid: true},
		IsActive:   true,
		Gender:     GenderMale,
		Generation: int32(util.RandomInt(2010, 2030)),
		Photo:      pgtype.Text{String: util.RandomString(12), Valid: true},
	}
	santri, err := testQueries.CreateSantri(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, santri)

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
		Name:       util.RandomString(8),
		Nis:        pgtype.Text{String: util.RandomString(15), Valid: true},
		IsActive:   true,
		Gender:     GenderMale,
		Generation: int32(util.RandomInt(2010, 2030)),
		Photo:      pgtype.Text{String: util.RandomString(12), Valid: true},
		ParentID:   pgtype.Int4{Int32: parent.ID, Valid: true},
	}
	santri, err := testQueries.CreateSantri(context.Background(), arg)
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
	clearSantriTable(t)
	createRandomSantri(t)
}

func TestListSantri(t *testing.T) {
	clearSantriTable(t)
	randomSantri, randomParent := createRandomSantriWithParent(t)
	santris := []Santri{}
	for i := 0; i < 10; i++ {
		santris = append(santris, createRandomSantri(t))
	}

	t.Run("Run with List name", func(t *testing.T) {
		arg := ListSantriAscNameParams{
			Q:            pgtype.Text{String: randomSantri.Name[:3], Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
		}

		allSantri, err := testQueries.ListSantriAscName(context.Background(), arg)
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

	t.Run("Run with List Nis", func(t *testing.T) {
		arg := ListSantriAscNameParams{
			Q:            pgtype.Text{String: randomSantri.Nis.String, Valid: true},
			LimitNumber:  10,
			OffsetNumber: 0,
		}

		allSantri, err := testQueries.ListSantriAscName(context.Background(), arg)
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

	t.Run("Run with List Parent Id", func(t *testing.T) {
		arg := ListSantriAscNameParams{
			ParentID:     randomSantri.ParentID,
			LimitNumber:  10,
			OffsetNumber: 0,
		}

		allSantri, err := testQueries.ListSantriAscName(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, allSantri)

		for _, santri := range allSantri {
			require.Equal(t, randomSantri.ParentID, santri.ParentID)
			require.NotZero(t, santri.ParentID)
			require.Equal(t, randomParent.Name, santri.ParentName.String)
		}
	})

	require.Equal(t, len(santris), 10)
}

func TestListSantriPagination(t *testing.T) {
	clearSantriTable(t)
	for i := 0; i < 15; i++ {
		createRandomSantri(t)
	}

	testCases := []struct {
		name     string
		arg      ListSantriAscNameParams
		expected int
	}{
		{
			name: "Limit 5",
			arg: ListSantriAscNameParams{
				LimitNumber:  5,
				OffsetNumber: 0,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 5",
			arg: ListSantriAscNameParams{
				LimitNumber:  5,
				OffsetNumber: 5,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 10",
			arg: ListSantriAscNameParams{
				LimitNumber:  5,
				OffsetNumber: 10,
			},
			expected: 5,
		},
		{
			name: "Limit 5 Offset 10",
			arg: ListSantriAscNameParams{
				LimitNumber:  5,
				OffsetNumber: 15,
			},
			expected: 0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			allSantri, err := testQueries.ListSantriAscName(context.Background(), tt.arg)
			require.NoError(t, err)
			require.Len(t, allSantri, tt.expected)
		})
	}
}

func TestUpdateSantri(t *testing.T) {
	clearSantriTable(t)
	santri := createRandomSantri(t)

	arg := UpdateSantriParams{
		ID:         santri.ID,
		Name:       util.RandomString(8),
		Nis:        pgtype.Text{String: util.RandomString(15), Valid: true},
		IsActive:   false,
		Generation: int32(util.RandomInt(2010, 2030)),
		Photo:      pgtype.Text{String: util.RandomString(12), Valid: true},
	}

	updatedSantri, err := testQueries.UpdateSantri(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSantri)

	require.Equal(t, arg.Name, updatedSantri.Name)
	require.Equal(t, arg.Nis.String, updatedSantri.Nis.String)
	require.Equal(t, arg.IsActive, updatedSantri.IsActive)
	require.Equal(t, arg.Generation, updatedSantri.Generation)
	require.Equal(t, arg.Photo.String, updatedSantri.Photo.String)
}

func TestDeleteSantri(t *testing.T) {
	clearSantriTable(t)
	santri := createRandomSantri(t)

	deletedSantri, err := testQueries.DeleteSantri(context.Background(), santri.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedSantri)

	getSantri, err := testQueries.GetSantri(context.Background(), santri.ID)
	require.Error(t, err)
	require.Empty(t, getSantri)

}
