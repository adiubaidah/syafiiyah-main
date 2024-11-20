package handler

import (
	"strconv"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/platform/mqtt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ArduinoHandler interface {
	CreateArduinoHandler(c *gin.Context)
	ListArduinosHandler(c *gin.Context)
	UpdateArduinoHandler(c *gin.Context)
	DeleteArduinoHandler(c *gin.Context)
}

type arduinoHandler struct {
	logger      *logrus.Logger
	usecase     usecase.ArduinoUseCase
	mqttHandler *mqtt.MQTTHandler
}

func NewArduinoHandler(logger *logrus.Logger, usecase usecase.ArduinoUseCase, mqttHandler *mqtt.MQTTHandler) ArduinoHandler {
	return &arduinoHandler{
		logger:      logger,
		usecase:     usecase,
		mqttHandler: mqttHandler,
	}
}

func (h *arduinoHandler) CreateArduinoHandler(c *gin.Context) {
	var request model.CreateArduinoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	arduino, err := h.usecase.CreateArduino(c, &request)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	h.mqttHandler.RefreshTopics()

	c.JSON(200, model.ResponseData[model.ArduinoResponse]{Code: 200, Status: "success", Data: arduino})
}

func (h *arduinoHandler) ListArduinosHandler(c *gin.Context) {
	arduinos, err := h.usecase.ListArduinos(c)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[[]model.ArduinoWithModesResponse]{Code: 200, Status: "success", Data: arduinos})
}

func (h *arduinoHandler) UpdateArduinoHandler(c *gin.Context) {
	idParam := c.Param("id")
	arduinoId, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	var request model.CreateArduinoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	arduino, err := h.usecase.UpdateArduino(c, &request, int32(arduinoId))
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	h.mqttHandler.RefreshTopics()

	c.JSON(200, model.ResponseData[model.ArduinoResponse]{Code: 200, Status: "success", Data: arduino})
}

func (h *arduinoHandler) DeleteArduinoHandler(c *gin.Context) {
	idParam := c.Param("id")
	arduinoId, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}
	arduino, err := h.usecase.DeleteArduino(c, int32(arduinoId))
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	h.mqttHandler.RefreshTopics()

	c.JSON(200, model.ResponseData[model.ArduinoResponse]{Code: 200, Status: "success", Data: arduino})
}
