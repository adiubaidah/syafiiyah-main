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
	// Handle 'presence' topic
}

func (h *MQTTHandler) handlePermission(client mqtt.Client, acknowledgmentTopic string, request *model.SmartCardRequest) {
	// Handle 'permission' topic
}

func (h *MQTTHandler) handlePing(client mqtt.Client, acknowledgmentTopic string) {
	// Handle 'ping' topic
}
