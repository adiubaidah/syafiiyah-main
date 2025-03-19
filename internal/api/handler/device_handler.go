package handler

import (
	"strconv"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	"github.com/adiubaidah/syafiiyah-main/internal/usecase"
	"github.com/adiubaidah/syafiiyah-main/platform/mqtt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DeviceHandler struct {
	Logger      *logrus.Logger
	UseCase     *usecase.DeviceUseCase
	MqttHandler *mqtt.MQTTBroker
}

func NewDeviceHandler(args *DeviceHandler) *DeviceHandler {
	return args
}

func (h *DeviceHandler) CreateDeviceHandler(c *gin.Context) {
	var request model.CreateDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	device, err := h.UseCase.CreateDevice(c, &request)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	h.MqttHandler.RefreshTopics()

	c.JSON(200, model.ResponseData[*model.DeviceResponse]{Code: 200, Status: "success", Data: device})
}

func (h *DeviceHandler) ListDevicesHandler(c *gin.Context) {
	result, err := h.UseCase.ListDevices(c)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	c.JSON(200, model.ResponseData[*[]model.DeviceWithModesResponse]{Code: 200, Status: "success", Data: result})
}

func (h *DeviceHandler) UpdateDeviceHandler(c *gin.Context) {
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

	device, err := h.UseCase.UpdateDevice(c, &request, int32(arduinoId))
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	h.MqttHandler.RefreshTopics()

	c.JSON(200, model.ResponseData[*model.DeviceResponse]{Code: 200, Status: "success", Data: device})
}

func (h *DeviceHandler) DeleteDeviceHandler(c *gin.Context) {
	idParam := c.Param("id")
	arduinoId, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}
	device, err := h.UseCase.DeleteDevice(c, int32(arduinoId))
	if err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	h.MqttHandler.RefreshTopics()

	c.JSON(200, model.ResponseData[*model.DeviceResponse]{Code: 200, Status: "success", Data: device})
}
