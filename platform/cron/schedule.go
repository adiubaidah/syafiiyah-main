package cron

import (
	"context"
	"errors"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/internal/worker"
	"github.com/sirupsen/logrus"
)

type ScheduleCron struct {
	logger                *logrus.Logger
	ActiveScheduleSantri  *model.SantriScheduleResponse
	santriScheduleUseCase usecase.SantriScheduleUseCase
	santriPresenceWorker  worker.SantriPresenceWorker
	stopChan              chan struct{}
	isRunning             bool
}

func NewScheduleCron(logger *logrus.Logger, santriScheduleUseCase usecase.SantriScheduleUseCase, santriPresenceWorker worker.SantriPresenceWorker) *ScheduleCron {
	return &ScheduleCron{
		santriScheduleUseCase: santriScheduleUseCase,
		santriPresenceWorker:  santriPresenceWorker,
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
	s.execute()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.execute()
		case <-s.stopChan:
			return
		}
	}
}

func (s *ScheduleCron) execute() {
	activeSchedule, err := s.santriScheduleUseCase.GetSantriSchedule(context.Background(), time.Now())
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			s.ActiveScheduleSantri = nil
		}
	}
	previousSchedule := s.ActiveScheduleSantri

	if !schedulesAreEqual(previousSchedule, activeSchedule) {
		s.logger.Infof("active schedule has changed: %v", activeSchedule)
		s.ActiveScheduleSantri = activeSchedule

		if previousSchedule != nil {
			go s.santriPresenceWorker.CreateAlphaForMissingPresence(previousSchedule.ID)
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
