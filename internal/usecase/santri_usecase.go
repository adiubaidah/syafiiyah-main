package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
)

type SantriUseCase struct {
	Log *logrus.Logger
}

func NewSantriUseCase(log *logrus.Logger) *SantriUseCase {
	return &SantriUseCase{Log: log}
}

func (c *SantriUseCase) Create(ctx context.Context) {

}
