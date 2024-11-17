package usecase

import (
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
)

type SantriPresenceUseCase interface {
}

type santriPresenceService struct {
	store db.Store
}

func NewSantriPresenceUseCase(store db.Store) SantriPresenceUseCase {
	return &santriPresenceService{
		store: store,
	}
}
