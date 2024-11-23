package mqtt

import (
	"context"
	"encoding/json"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (h *MQTTHandler) handleRecord(client mqtt.Client, acknowledgmentTopic string, request *model.SmartCardRequest) {
	recordedSmartCard, err := h.SmartCardUseCase.CreateSmartCard(context.Background(), request)
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
			Data:   recordedSmartCard,
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

	getSmartCard, err := h.SmartCardUseCase.GetSmartCard(context.Background(), &model.SmartCardRequest{Uid: request.Uid})
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

	if getSmartCard.OwnerRole == "" {
		h.logger.Warn("Smart card owner role is not set")

		response := model.ResponseMessage{
			Code:    403,
			Status:  "error",
			Message: "Owner belum ada",
		}
		h.publishResponse(client, acknowledgmentTopic, response)
		return
	} else {

	}

	h.logger.Infof("%s %s is present", getSmartCard.OwnerRole, getSmartCard.Details.Name)

	response := model.ResponseData[model.UserResponse]{
		Code:   200,
		Status: "success",
		Data: model.UserResponse{
			ID:       getSmartCard.Details.ID,
			Username: getSmartCard.Details.Name,
			Role:     getSmartCard.OwnerRole,
		},
	}
	h.publishResponse(client, acknowledgmentTopic, response)
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
