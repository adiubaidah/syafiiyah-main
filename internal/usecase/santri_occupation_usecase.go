package usecase

import (
	"context"
	"fmt"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/jackc/pgx/v5/pgtype"
)

type SantriOccuapationUsecase interface {
	CreateSantriOccupation(ctx context.Context, santriOccupation model.CreateSantriOccupationRequest) (model.SantriOccupationResponse, error)
}

type service struct {
	store db.Store
}

func NewSantriOccupationUseCase(store db.Store) SantriOccuapationUsecase {
	return &service{store: store}
}

func (s *service) CreateSantriOccupation(ctx context.Context, santriOccupation model.CreateSantriOccupationRequest) (model.SantriOccupationResponse, error) {

	fmt.Println("usecase")
	createdSantriOccupation, err := s.store.CreateSantriOccupation(ctx, db.CreateSantriOccupationParams{
		Name:        santriOccupation.Name,
		Description: pgtype.Text{String: santriOccupation.Description, Valid: true},
	})
	if err != nil {
		return model.SantriOccupationResponse{}, err
	}

	return model.SantriOccupationResponse{
		ID: createdSantriOccupation.ID,
		BaseSantriOccupation: model.BaseSantriOccupation{
			Name:        createdSantriOccupation.Name,
			Description: createdSantriOccupation.Description.String,
		},
	}, nil
}
