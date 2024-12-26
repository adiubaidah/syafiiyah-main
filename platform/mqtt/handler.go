package mqtt

import (
	"context"
	"encoding/json"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	repo "github.com/adiubaidah/rfid-syafiiyah/internal/repository"
)

func (h *MQTTBroker) handleRecord(acknowledgmentTopic string, request *model.SmartCardRequest) {
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

	token := h.Client.Publish(acknowledgmentTopic, 0, false, payload)

	if token.Wait() && token.Error() != nil {
		h.logger.Errorf("Error sending acknowledgment: %v\n", token.Error())
	}
}

func (h *MQTTBroker) handlePresence(acknowledgmentTopic string, request *model.SmartCardRequest) {

	getSmartCard, err := h.smartCardUseCase.GetSmartCard(context.Background(), &model.SmartCardRequest{Uid: request.Uid})
	var response any
	if err != nil {
		h.logger.Errorf("Error getting smart card: %v\n", err)

		response = createErrorResponse(err)
		h.publishResponse(acknowledgmentTopic, response)
		return
	}

	if !getSmartCard.IsActive {
		h.logger.Warn("Smart card is not active")

		response = model.ResponseMessage{
			Code:    403,
			Status:  "error",
			Message: "Smart card tidak aktif",
		}
		h.publishResponse(acknowledgmentTopic, response)
		return
	}

	switch getSmartCard.Owner.Role {
	case repo.RoleTypeSantri:
		result, err := h.SantriHandler.Presence(request.Uid, getSmartCard.Owner.ID)
		if err != nil {
			if appErr, ok := err.(*exception.AppError); ok {
				response = model.ResponseMessage{
					Code:    appErr.Code,
					Status:  "error",
					Message: appErr.Message,
				}
			} else {
				response = model.ResponseMessage{
					Code:    500,
					Status:  "error",
					Message: err.Error(),
				}
			}
		} else {
			response = model.ResponseData[*model.SantriPresenceResponse]{
				Code:   200,
				Status: "success",
				Data:   result,
			}
		}
		h.publishResponse(acknowledgmentTopic, response)
		return
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

func (h *MQTTBroker) publishResponse(topic string, response any) {
	payload, err := json.Marshal(response)
	if err != nil {
		h.logger.Errorf("Error marshaling response: %v\n", err)
		return
	}

	token := h.Client.Publish(topic, 0, false, payload)
	if token.Wait() && token.Error() != nil {
		h.logger.Errorf("Error sending acknowledgment: %v\n", token.Error())
	}
}

func (h *MQTTBroker) handlePermission(acknowledgmentTopic string, request *model.SmartCardRequest) {
	// Handle 'permission' topic
}

func (h *MQTTBroker) handlePing(acknowledgmentTopic string) {
	// Handle 'ping' topic
}
