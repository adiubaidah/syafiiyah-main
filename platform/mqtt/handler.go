package mqtt

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (h *MQTTHandler) handleRecord(client mqtt.Client, acknowledgmentTopic string, request *model.SmartCardRequest) {
	recordedSmartCard, err := h.smartCardUseCase.CreateSmartCard(context.Background(), request)
	if err != nil {
		h.logger.Errorf("Error creating smart card: %v\n", err)
	}

	var response any
	if err != nil {
		h.logger.Errorf("Error creating smart card: %v\n", err)
		if appError, ok := err.(*exception.AppError); ok {
			response = model.ResponseMessage{
				Code:    appError.Code,
				Status:  "error",
				Message: appError.Message,
			}
		} else {
			response = model.ResponseMessage{
				Code:    500,
				Status:  "error",
				Message: err.Error(),
			}
		}
	} else {

		response = model.ResponseData[model.SmartCard]{
			Code:   200,
			Status: "success",
			Data:   *recordedSmartCard,
		}
	}

	payload, err := json.Marshal(response)
	if err != nil {
		h.logger.Errorf("Error marshaling response: %v\n", err)
		return
	}

	token := client.Publish(acknowledgmentTopic, 0, false, payload)

	if token.Wait() && token.Error() != nil {
		h.logger.Errorf("Error sending acknowledgment: %v\n", token.Error())
	}
}

func (h *MQTTHandler) handlePresence(client mqtt.Client, acknowledgmentTopic string, request *model.SmartCardRequest) {

	CURRENT_TIME_PRESENCE := time.Now()

	h.logger.Println("Current time: ", CURRENT_TIME_PRESENCE)

	getSmartCard, err := h.smartCardUseCase.GetSmartCard(context.Background(), &model.SmartCardRequest{Uid: request.Uid})
	if err != nil {
		h.logger.Errorf("Error getting smart card: %v\n", err)

		response := createErrorResponse(err)
		h.publishResponse(client, acknowledgmentTopic, response)
		return
	}

	if !getSmartCard.IsActive {
		h.logger.Warn("Smart card is not active")

		response := model.ResponseMessage{
			Code:    403,
			Status:  "error",
			Message: "Smart card tidak aktif",
		}
		h.publishResponse(client, acknowledgmentTopic, response)
		return
	}

	if getSmartCard.Owner.Role == "" {
		h.logger.Warn("Smart card owner role is not set")

		response := model.ResponseMessage{
			Code:    403,
			Status:  "error",
			Message: "Owner belum ada",
		}
		h.publishResponse(client, acknowledgmentTopic, response)
		return
	} else {

		if getSmartCard.Owner.Role == db.RoleTypeSantri {
			// Check if the smart card is registered as a santri
			_, err := h.santriUseCase.GetSantri(context.Background(), getSmartCard.Owner.ID)
			if err != nil {
				h.logger.Errorf("Error getting santri: %v\n", err)

				response := createErrorResponse(err)
				h.publishResponse(client, acknowledgmentTopic, response)
				return
			}

			// check santri late or not according active santri schedule
			santriStartPresence, err := util.ParseTimeWithCurrentDate(h.schedule.ActiveScheduleSantri.StartPresence)
			if err != nil {
				h.logger.Errorf("Error parsing time: %v\n", err)
			}
			santriStartTime, _ := util.ParseTimeWithCurrentDate(h.schedule.ActiveScheduleSantri.StartTime)
			// santriFinishTime, _ := util.ParseTime(h.schedule.ActiveScheduleSantri.FinishTime)

			arg := &model.CreateSantriPresenceRequest{
				ScheduleID:   h.schedule.ActiveScheduleSantri.ID,
				ScheduleName: h.schedule.ActiveScheduleSantri.Name,
				SantriID:     getSmartCard.Owner.ID,
				CreatedBy:    db.PresenceCreatedByTypeTap,
				CreatedAt:    CURRENT_TIME_PRESENCE.String(),
			}
			if CURRENT_TIME_PRESENCE.After(santriStartPresence) && CURRENT_TIME_PRESENCE.Before(santriStartTime) {
				arg.Type = db.PresenceTypePresent
			} else if CURRENT_TIME_PRESENCE.After(santriStartTime) {
				arg.Type = db.PresenceTypeLate
			}
			presence, err := h.santriPresenceUseCase.CreateSantriPresence(context.Background(), arg)
			if err != nil {
				h.logger.Errorf("Error creating santri presence: %v\n", err)

				response := createErrorResponse(err)
				h.publishResponse(client, acknowledgmentTopic, response)
				return
			}

			response := model.ResponseData[model.SantriPresenceResponse]{
				Code:   200,
				Status: "success",
				Data:   *presence,
			}
			h.publishResponse(client, acknowledgmentTopic, response)
		} else {
			h.logger.Warnf("Unknown owner role: %s", getSmartCard.Owner.Role)

			response := model.ResponseMessage{
				Code:    403,
				Status:  "error",
				Message: "Owner tidak dikenal",
			}
			h.publishResponse(client, acknowledgmentTopic, response)
			return
		}
	}
}

func createErrorResponse(err error) model.ResponseMessage {
	if appErr, ok := err.(*exception.AppError); ok {
		return model.ResponseMessage{
			Code:    appErr.Code,
			Status:  "error",
			Message: appErr.Message,
		}
	}

	return model.ResponseMessage{
		Code:    500,
		Status:  "error",
		Message: err.Error(),
	}
}

func (h *MQTTHandler) publishResponse(client mqtt.Client, topic string, response any) {
	payload, err := json.Marshal(response)
	if err != nil {
		h.logger.Errorf("Error marshaling response: %v\n", err)
		return
	}

	token := client.Publish(topic, 0, false, payload)
	if token.Wait() && token.Error() != nil {
		h.logger.Errorf("Error sending acknowledgment: %v\n", token.Error())
	}
}

func (h *MQTTHandler) handlePermission(client mqtt.Client, acknowledgmentTopic string, request *model.SmartCardRequest) {
	// Handle 'permission' topic
}

func (h *MQTTHandler) handlePing(client mqtt.Client, acknowledgmentTopic string) {
	// Handle 'ping' topic
}
