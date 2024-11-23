package cron

import (
	"context"
	"errors"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/sirupsen/logrus"
)

type ScheduleCron struct {
	logger                *logrus.Logger
	ActiveScheduleSantri  model.SantriScheduleResponse
	santriScheduleUseCase usecase.SantriScheduleUseCase
	stopChan              chan struct{}
	isRunning             bool
}

func NewScheduleCron(logger *logrus.Logger, santriScheduleUseCase usecase.SantriScheduleUseCase) *ScheduleCron {
	return &ScheduleCron{
		santriScheduleUseCase: santriScheduleUseCase,
		logger:                logger,
	}
}

func (s *ScheduleCron) Start() {
	if s.isRunning {
		return
	}
	s.stopChan = make(chan struct{})
	s.isRunning = true
	go s.run()
}

func (s *ScheduleCron) Stop() {
	if !s.isRunning {
		return
	}
	close(s.stopChan)
	s.isRunning = false
}

func (s *ScheduleCron) run() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Get active schedule
			activeSchedule, err := s.santriScheduleUseCase.GetSantriSchedule(context.Background(), time.Now())
			s.logger.Println(time.Now())
			if err != nil {
				s.logger.Errorf("failed to get active schedule: %v", err)
				if errors.Is(err, exception.ErrNotFound) {
					s.ActiveScheduleSantri = model.SantriScheduleResponse{}
				}
			} else {
				s.ActiveScheduleSantri = activeSchedule
			}
		case <-s.stopChan:
			return
		}
	}
}
