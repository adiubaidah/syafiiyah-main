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
	ActiveScheduleSantri  *model.SantriScheduleResponse
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

	// Wait until the next minute
	now := time.Now()
	nextMinute := now.Truncate(time.Minute).Add(time.Minute)
	timeUntilNextMinute := time.Until(nextMinute)
	time.Sleep(timeUntilNextMinute)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Get active schedule
			activeSchedule, err := s.santriScheduleUseCase.GetSantriSchedule(context.Background(), time.Now())
			if err != nil {
				s.logger.Errorf("failed to get active schedule: %v", err)
				if errors.Is(err, exception.ErrNotFound) {
					s.ActiveScheduleSantri = nil
				}
			}
			previousSchedule := s.ActiveScheduleSantri

			if !schedulesAreEqual(previousSchedule, activeSchedule) {
				s.logger.Infof("active schedule has changed: %v", activeSchedule)
				s.ActiveScheduleSantri = activeSchedule

				// Notify schedule change
				if previousSchedule != nil {

				}

			}

		case <-s.stopChan:
			return
		}
	}
}

func schedulesAreEqual(a, b *model.SantriScheduleResponse) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.ID == b.ID && a.StartPresence == b.StartPresence
}
