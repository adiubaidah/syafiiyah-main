package usecase

import (
	"context"
	"testing"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	mocks "github.com/adiubaidah/syafiiyah-main/internal/repository/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEmployeeUseCase_Create(t *testing.T) {
	mockStore := new(mocks.MockStore)
	uc := NewEmployeeUseCase(mockStore)
	ctx := context.Background()

	request := &model.CreateEmployeeRequest{
		NIP:          "12345",
		Name:         "John Doe",
		Gender:       "Male",
		Photo:        "photo.jpg",
		OccupationID: 1,
		UserID:       2,
	}

	expectedResult := repo.Employee{
		ID:           1,
		Nip:          pgtype.Text{String: "12345", Valid: true},
		Name:         "John Doe",
		Gender:       "Male",
		Photo:        pgtype.Text{String: "photo.jpg", Valid: true},
		OccupationID: 1,
		UserID:       pgtype.Int4{Int32: 2, Valid: true},
	}

	// Success Scenario
	mockStore.On("CreateEmployee", ctx, mock.Anything).Return(expectedResult, nil)

	result, err := uc.Create(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedResult.ID, result.ID)
	require.Equal(t, expectedResult.Nip.String, result.NIP)
	require.Equal(t, expectedResult.Name, result.Name)
}
