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

type EmployeeMQTTHandler struct {
	logger          *logrus.Logger
	usecase         *usecase.EmployeeUseCase
	presenceUseCase usecase.EmployeePresenceUseCase
	service         pb.EmployeeScheduleServiceClient
}

func NewEmployeeMQTTHandler(arg *EmployeeMQTTHandler) *EmployeeMQTTHandler {
	return arg
}

func (h *EmployeeMQTTHandler) Presence(uid string, santriID int32) (*model.EmployeePresenceResponse, error) {

	activeSchedule, err := h.service.ActiveEmployeeSchedule(context.Background(), &pb.ActiveEmployeeScheduleRequest{})
	if err != nil {
		return nil, exception.NewNotFoundError("no active schedule found for employee attendance")
	}

	CURRENT_TIME_PRESENCE := time.Now()
	_, err = h.usecase.GetByID(context.Background(), santriID)
	if err != nil {
		h.logger.Errorf("Error getting employee: %v\n", err)
	}

	startPresence, err := util.ParseHHMMWithCurrentDate(activeSchedule.StartPresence)
	if err != nil {
		h.logger.Errorf("Error parsing time: %v\n", err)
		return nil, exception.NewParseTimeError("start presence", err)
	}
	startTime, _ := util.ParseHHMMWithCurrentDate(activeSchedule.StartTime)

	arg := &model.CreateEmployeePresenceRequest{
		ScheduleID:   activeSchedule.Id,
		ScheduleName: activeSchedule.Name,
		EmployeeID:   santriID,
		CreatedBy:    repo.PresenceCreatedByTypeTap,
	}
	if CURRENT_TIME_PRESENCE.After(startPresence) && CURRENT_TIME_PRESENCE.Before(startTime) {
		arg.Type = repo.PresenceTypePresent
	} else if CURRENT_TIME_PRESENCE.After(startTime) {
		arg.Type = repo.PresenceTypeLate
	}
	presence, err := h.presenceUseCase.CreatePresence(context.Background(), arg)
	if err != nil {
		h.logger.Errorf("Error creating employee presence: %v\n", err)
		if exception.DatabaseErrorCode(err) == exception.ErrCodeUniqueViolation {
			return nil, exception.NewUniqueViolationError("employee already presence today", err)
		}
	}

	return presence, nil
}
