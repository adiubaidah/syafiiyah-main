package handler

import (
	"strconv"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	"github.com/adiubaidah/rfid-syafiiyah/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HolidayHandler interface {
	CreateHolidayHandler(c *gin.Context)
	ListHolidaysHandler(c *gin.Context)
	UpdateHolidayHandler(c *gin.Context)
	DeleteHolidayHandler(c *gin.Context)
}

type holidayHandler struct {
	logger  *logrus.Logger
	usecase usecase.HolidayUseCase
}

func NewHolidayHandler(logger *logrus.Logger, usecase usecase.HolidayUseCase) HolidayHandler {
	return &holidayHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (h *holidayHandler) CreateHolidayHandler(c *gin.Context) {
	var request model.CreateHolidayRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	holiday, err := h.usecase.CreateHoliday(c, &request)
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	c.JSON(201, model.ResponseData[*model.HolidayResponse]{Code: 201, Status: "success", Data: holiday})
}

func (h *holidayHandler) ListHolidaysHandler(c *gin.Context) {
	var request model.ListHolidayRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	holidays, err := h.usecase.ListHolidays(c, &request)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
	}

	c.JSON(200, model.ResponseData[*[]model.HolidayResponse]{Code: 200, Status: "success", Data: holidays})
}

func (h *holidayHandler) UpdateHolidayHandler(c *gin.Context) {
	var request model.UpdateHolidayRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}

	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}
	holidayId := int32(idParam)
	updatedHoliday, err := h.usecase.UpdateHoliday(c, &request, holidayId)
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			h.logger.Error(appErr.Message)
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	c.JSON(200, model.ResponseData[*model.HolidayResponse]{Code: 200, Status: "success", Data: updatedHoliday})
}

func (h *holidayHandler) DeleteHolidayHandler(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, model.ResponseMessage{Code: 400, Status: "error", Message: err.Error()})
		return
	}
	holidayId := int32(idParam)
	deletedHoliday, err := h.usecase.DeleteHoliday(c, holidayId)
	if err != nil {
		h.logger.Error(err)
		if appErr, ok := err.(*exception.AppError); ok {
			h.logger.Error(appErr.Message)
			c.JSON(appErr.Code, model.ResponseMessage{Code: appErr.Code, Status: "error", Message: appErr.Message})
			return
		}
		c.JSON(500, model.ResponseMessage{Code: 500, Status: "error", Message: "Internal server error"})
		return
	}

	c.JSON(200, model.ResponseData[*model.HolidayResponse]{Code: 200, Status: "sucess", Data: deletedHoliday})
}
