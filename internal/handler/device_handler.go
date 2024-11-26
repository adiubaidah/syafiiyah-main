package handler

import (
	"strconv"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/adiubaidah/rfid-syafiiyah/platform/mqtt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DeviceHandler interface {
	CreateDeviceHandler(c *gin.Context)
	ListDevicesHandler(c *gin.Context)
	UpdateDeviceHandler(c *gin.Context)
	DeleteDeviceHandler(c *gin.Context)
}

type deviceHandler struct {
	logger      *logrus.Logger
	usecase     usecase.DeviceUseCase
	mqttHandler *mqtt.MQTTHandler
}

func NewDeviceHandler(logger *logrus.Logger, usecase usecase.DeviceUseCase, mqttHandler *mqtt.MQTTHandler) DeviceHandler {
	return &deviceHandler{
		logger:      logger,
		usecase:     usecase,
		mqttHandler: mqttHandler,
	}
}

func (h *deviceHandler) CreateDeviceHandler(c *gin.Context) {
	var request model.CreateDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	device, err := h.usecase.CreateDevice(c, &request)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	h.mqttHandler.RefreshTopics()

	c.JSON(200, model.ResponseData[model.DeviceResponse]{Code: 200, Status: "success", Data: device})
}

func (h *deviceHandler) ListDevicesHandler(c *gin.Context) {
	result, err := h.usecase.ListDevices(c)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[[]model.DeviceWithModesResponse]{Code: 200, Status: "success", Data: result})
}

func (h *deviceHandler) UpdateDeviceHandler(c *gin.Context) {
	idParam := c.Param("id")
	arduinoId, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	var request model.CreateDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	arduino, err := h.usecase.UpdateDevice(c, &request, int32(arduinoId))
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	h.mqttHandler.RefreshTopics()

	c.JSON(200, model.ResponseData[model.DeviceResponse]{Code: 200, Status: "success", Data: arduino})
}

func (h *deviceHandler) DeleteDeviceHandler(c *gin.Context) {
	idParam := c.Param("id")
	arduinoId, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}
	arduino, err := h.usecase.DeleteDevice(c, int32(arduinoId))
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	h.mqttHandler.RefreshTopics()

	c.JSON(200, model.ResponseData[model.DeviceResponse]{Code: 200, Status: "success", Data: arduino})
}
