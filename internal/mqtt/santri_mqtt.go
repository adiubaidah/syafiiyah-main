package mqtt_handler

import (
	"context"
	"time"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/exception"
	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	pb "github.com/adiubaidah/syafiiyah-main/internal/protobuf"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/adiubaidah/syafiiyah-main/internal/usecase"
	"github.com/adiubaidah/syafiiyah-main/pkg/util"
	"github.com/sirupsen/logrus"
)

type SantriMQTTHandler struct {
	logger          *logrus.Logger
	usecase         usecase.SantriUseCase
	presenceUseCase usecase.SantriPresenceUseCase
	service         pb.SantriScheduleServiceClient
}

func NewSantriMQTTHandler(logger *logrus.Logger, usecase usecase.SantriUseCase, service pb.SantriScheduleServiceClient, presenceUseCase usecase.SantriPresenceUseCase) *SantriMQTTHandler {
	return &SantriMQTTHandler{
		logger:          logger,
		usecase:         usecase,
		presenceUseCase: presenceUseCase,
		service:         service,
	}
}

func (h *SantriMQTTHandler) Presence(uid string, santriID int32) (*model.SantriPresenceResponse, error) {

	activeSchedule, err := h.service.ActiveSantriSchedule(context.Background(), &pb.ActiveSantriScheduleRequest{})
	if err != nil {
		return nil, exception.NewNotFoundError("no active schedule found for santri attendance")
	}

	CURRENT_TIME_PRESENCE := time.Now()
	_, err = h.usecase.GetSantri(context.Background(), santriID)
	if err != nil {
		h.logger.Errorf("Error getting santri: %v\n", err)
	}

	santriStartPresence, err := util.ParseHHMMWithCurrentDate(activeSchedule.StartPresence)
	if err != nil {
		h.logger.Errorf("Error parsing time: %v\n", err)
		return nil, exception.NewParseTimeError("start presence", err)
	}
	santriStartTime, _ := util.ParseHHMMWithCurrentDate(activeSchedule.StartTime)

	arg := &model.CreateSantriPresenceRequest{
		ScheduleID:   activeSchedule.Id,
		ScheduleName: activeSchedule.Name,
		SantriID:     santriID,
		CreatedBy:    repo.PresenceCreatedByTypeTap,
	}
	if CURRENT_TIME_PRESENCE.After(santriStartPresence) && CURRENT_TIME_PRESENCE.Before(santriStartTime) {
		arg.Type = repo.PresenceTypePresent
	} else if CURRENT_TIME_PRESENCE.After(santriStartTime) {
		arg.Type = repo.PresenceTypeLate
	}
	presence, err := h.presenceUseCase.CreateSantriPresence(context.Background(), arg)
	if err != nil {
		h.logger.Errorf("Error creating santri presence: %v\n", err)
		if exception.DatabaseErrorCode(err) == exception.ErrCodeUniqueViolation {
			return nil, exception.NewUniqueViolationError("santri already presence today", err)
		}
	}

	return presence, nil
}
