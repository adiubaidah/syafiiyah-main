package mqtt_handler

import (
	"context"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/adiubaidah/rfid-syafiiyah/platform/cron"
	"github.com/sirupsen/logrus"
)

type SantriMQTTHandler struct {
	logger          *logrus.Logger
	usecase         usecase.SantriUseCase
	presenceUseCase usecase.SantriPresenceUseCase
	schedule        *cron.ScheduleCron
}

func NewSantriMQTTHandler(logger *logrus.Logger, schedule *cron.ScheduleCron, usecase usecase.SantriUseCase, presenceUseCase usecase.SantriPresenceUseCase) *SantriMQTTHandler {
	return &SantriMQTTHandler{
		logger:          logger,
		usecase:         usecase,
		presenceUseCase: presenceUseCase,
		schedule:        schedule,
	}
}

func (h *SantriMQTTHandler) Presence(uid string, santriID int32) (*model.SantriPresenceResponse, error) {
	if h.schedule.ActiveScheduleSantri == nil {
		return nil, exception.NewNotFoundError("no active schedule found for santri attendance")
	}

	CURRENT_TIME_PRESENCE := time.Now()
	_, err := h.usecase.GetSantri(context.Background(), santriID)
	if err != nil {
		h.logger.Errorf("Error getting santri: %v\n", err)
	}

	santriStartPresence, err := util.ParseTimeWithCurrentDate(h.schedule.ActiveScheduleSantri.StartPresence)
	if err != nil {
		h.logger.Errorf("Error parsing time: %v\n", err)
		return nil, exception.NewParseTimeError("start presence", err)
	}
	santriStartTime, _ := util.ParseTimeWithCurrentDate(h.schedule.ActiveScheduleSantri.StartTime)

	arg := &model.CreateSantriPresenceRequest{
		ScheduleID:   h.schedule.ActiveScheduleSantri.ID,
		ScheduleName: h.schedule.ActiveScheduleSantri.Name,
		SantriID:     santriID,
		CreatedBy:    db.PresenceCreatedByTypeTap,
	}
	if CURRENT_TIME_PRESENCE.After(santriStartPresence) && CURRENT_TIME_PRESENCE.Before(santriStartTime) {
		arg.Type = db.PresenceTypePresent
	} else if CURRENT_TIME_PRESENCE.After(santriStartTime) {
		arg.Type = db.PresenceTypeLate
	}
	presence, err := h.presenceUseCase.CreateSantriPresence(context.Background(), arg)
	if err != nil {
		h.logger.Errorf("Error creating santri presence: %v\n", err)
		if exception.DatabaseErrorCode(err) == exception.ErrCodeUniqueViolation {
			return nil, exception.NewUniqueViolationError("err", err)
		}
	}

	return presence, nil
}
